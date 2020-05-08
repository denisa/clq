package validator

import (
	"errors"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A Config struct has configurations for the Validator renderer.
type Config struct {
	Release bool
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{
		Release: false,
	}
}

// SetOption implements renderer.NodeRenderer.SetOption.
func (c *Config) SetOption(name renderer.OptionName, value interface{}) {
	switch name {
	case optRelease:
		c.Release = value.(bool)
	}
}

// An Option interface sets options for Validator based renderers.
type Option interface {
	SetValidationOption(*Config)
}

// Release is an option that control the validation mode
const optRelease renderer.OptionName = "Release"

type withRelease struct {
	value bool
}

func (o *withRelease) SetConfig(c *renderer.Config) {
	c.Options[optRelease] = o.value
}

func (o *withRelease) SetValidationOption(c *Config) {
	c.Release = o.value
}

// WithRelease is a functional option that allow you to set the release mode to
// the renderer.
func WithRelease(release bool) interface {
	renderer.Option
	Option
} {
	return &withRelease{release}
}

// A Renderer struct is an implementation of renderer.NodeRenderer that validates
// a changelog.
type Renderer struct {
	Config
	text                     strings.Builder
	h1Released, h1Unreleased bool
	changes                  changelog.ChangeMap
	hasChangeDescriptions    bool
	headers                  changelog.Stack
	previousRelease          changelog.Release
}

func NewRenderer(opts ...Option) renderer.NodeRenderer {
	r := &Renderer{
		Config:  NewConfig(),
		changes: make(changelog.ChangeMap),
		headers: changelog.NewStack(),
	}

	for _, opt := range opts {
		opt.SetValidationOption(&r.Config)
	}
	return r
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDocument, r.visitDocument)
	reg.Register(ast.KindHeading, r.visitHeading)
	reg.Register(ast.KindList, r.visitList)
	reg.Register(ast.KindListItem, r.visitListItem)

	reg.Register(ast.KindImage, r.visitImage)
	reg.Register(ast.KindText, r.visitText)
}

func (r *Renderer) visitDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if !r.h1Released && !r.h1Unreleased {
			return ast.WalkStop, errors.New("Validation error: No release defined in changelog")
		}
		if (r.headers.Release() || r.headers.Change()) && !r.hasChangeDescriptions {
			return ast.WalkStop, errors.New("No change descriptions for " + r.headers.AsPath())
		}
		w.WriteString(r.queryEngine.Result())
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.text.Reset()
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Heading)
	switch n.Level {
	case 1:
		h, err := r.headers.ResetTo(changelog.TitleHeading, r.text.String())
		if err != nil {
			return ast.WalkStop, err
		}

		title := h.(changelog.Changelog)
		if err := r.validateDocumentHeading(title); err != nil {
			return ast.WalkStop, err
		}
	case 2:
		if (r.headers.Release() || r.headers.Change()) && !r.hasChangeDescriptions {
			if err := errors.New("No change descriptions for " + r.headers.AsPath()); err != nil {
				return ast.WalkStop, err
			}
		} else {
			h, err := r.headers.ResetTo(changelog.ReleaseHeading, r.text.String())
			if err != nil {
				return ast.WalkStop, err
			}

			release := h.(changelog.Release)
			if err := r.validateReleaseHeading(release); err != nil {
				return ast.WalkStop, err
			}

			if r.previousRelease.IsRelease() && release.IsRelease() {
				nextRelease := release.NextRelease(r.changes)
				if !r.previousRelease.HasRelease(nextRelease) {
					return ast.WalkStop, errors.New("Release '" + r.previousRelease.Name() + "' should have version " + nextRelease.String())
				}
			}

			if release.Unreleased() {
				r.h1Unreleased = true
			} else if !release.Yanked() {
				r.h1Released = true
			}
			r.hasChangeDescriptions = false
			r.changes = make(changelog.ChangeMap)
			r.previousRelease = release
		}
	case 3:
		if r.headers.Title() {
			return ast.WalkStop, errors.New("Changes must be in a release " + r.headers.AsPath())
		} else if r.headers.Change() && !r.hasChangeDescriptions {
			return ast.WalkStop, errors.New("No change descriptions for " + r.headers.AsPath())
		} else {
			h, err := r.headers.ResetTo(changelog.ChangeHeading, r.text.String())
			if err != nil {
				return ast.WalkStop, err
			}

			change := h.(changelog.Change)
			if err := r.validateChangeHeading(change); err != nil {
				return ast.WalkStop, err
			}
			r.hasChangeDescriptions = false
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) validateDocumentHeading(title changelog.Changelog) error { return nil }

func (r *Renderer) validateReleaseHeading(release changelog.Release) error {
	if release.Unreleased() {
		if r.Release {
			return errors.New("Validation error: [Unreleased] not supported in release mode " + r.headers.AsPath())
		}
		if r.h1Unreleased {
			return errors.New("Validation error: Multiple [Unreleased] not supported " + r.headers.AsPath())
		}
		if r.h1Released {
			return errors.New("Validation error: [Unreleased] must come before any release " + r.headers.AsPath())
		}
	} else {
		if release.Yanked() {
			if !r.h1Released && !r.h1Unreleased {
				return errors.New("Validation error: Changelog cannot start with a [YANKED] release, insert a release or a [Unreleased] first " + r.headers.AsPath())
			}
		}
		if r.previousRelease.HasBeenReleased() {
			if err := r.previousRelease.SortsBefore(release); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Renderer) validateChangeHeading(change changelog.Change) error {
	if r.changes[change.Name()] {
		return errors.New("Validation error: Multiple headings " + change.Name() + " not supported " + r.headers.AsPath())
	}
	r.changes[change.Name()] = true
	return nil
}

func (r *Renderer) visitList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if r.headers.Change() {
		r.hasChangeDescriptions = true
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Text)
		segment := n.Segment
		r.text.Write(segment.Value(source))
	}
	return ast.WalkContinue, nil
}
