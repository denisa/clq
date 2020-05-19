package validator

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/query"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A Config struct has configurations for the Validator renderer.
type Config struct {
	release     bool
	queryEngine query.QueryEngine
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{}
}

// An Option interface sets options for Validator based renderers.
type Option interface {
	SetValidationOption(*Config)
}

// Query is the optional query
const optQuery renderer.OptionName = "Query"

type withQuery struct {
	value query.QueryEngine
}

func (o *withQuery) SetValidationOption(c *Config) {
	c.queryEngine = o.value
}

// WithQuery is a functional option that allow you to set the query string to
// the renderer.
func WithQuery(queryEngine query.QueryEngine) interface {
	Option
} {
	return &withQuery{queryEngine}
}

// Release is an option that control the validation mode
const optRelease renderer.OptionName = "Release"

type withRelease struct {
	value bool
}

func (o *withRelease) SetValidationOption(c *Config) {
	c.release = o.value
}

// WithRelease is a functional option that allow you to set the release mode to
// the renderer.
func WithRelease(release bool) interface {
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
			return ast.WalkStop, fmt.Errorf("Validation error: No release defined in changelog")
		}
		if (r.headers.Release() || r.headers.Change()) && !r.hasChangeDescriptions {
			return ast.WalkStop, fmt.Errorf("No change descriptions for %v", r.headers)
		}
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
		// no validation rules defined for the title...

		r.queryEngine.Apply(w, title)
	case 2:
		if (r.headers.Release() || r.headers.Change()) && !r.hasChangeDescriptions {
			if err := fmt.Errorf("No change descriptions for %v", r.headers); err != nil {
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
					return ast.WalkStop, fmt.Errorf("Release %q should have version %v", r.previousRelease.Name(), nextRelease)
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

			r.queryEngine.Apply(w, release)
		}
	case 3:
		if r.headers.Title() {
			return ast.WalkStop, fmt.Errorf("Changes must be in a release %v", r.headers)
		} else if r.headers.Change() && !r.hasChangeDescriptions {
			return ast.WalkStop, fmt.Errorf("No change descriptions for %v", r.headers)
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

			r.queryEngine.Apply(w, change)
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) validateReleaseHeading(release changelog.Release) error {
	if release.Unreleased() {
		if r.release {
			return fmt.Errorf("Validation error: \"[Unreleased]\" not supported in release mode %v", r.headers)
		}
		if r.h1Unreleased {
			return fmt.Errorf("Validation error: Multiple \"[Unreleased]\" not supported %v", r.headers)
		}
		if r.h1Released {
			return fmt.Errorf("Validation error: \"[Unreleased]\" must come before any release %v", r.headers)
		}
	} else {
		if release.Yanked() {
			if !r.h1Released && !r.h1Unreleased {
				return fmt.Errorf("Validation error: Changelog cannot start with a \"[YANKED]\" release, insert a release or a \"[Unreleased]\" first %v", r.headers)
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
		return fmt.Errorf("Validation error: Multiple headings %q not supported %v", change.Name(), r.headers)
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
