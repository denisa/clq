package validator

import (
	"errors"
	"strings"
	"time"

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

type changeMap map[string]bool

// A Renderer struct is an implementation of renderer.NodeRenderer that validates
// a changelog.
type Renderer struct {
	Config
	text                     strings.Builder
	h1Released, h1Unreleased bool
	changes                  changeMap
	hasChangeDescriptions    bool
	headers                  stack
	previousRelease          Release
}

func NewRenderer(opts ...Option) renderer.NodeRenderer {
	r := &Renderer{
		Config:  NewConfig(),
		changes: make(changeMap),
		headers: NewStack(),
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
	var err error
	if !entering {
		if !r.h1Released && !r.h1Unreleased {
			err = errors.New("Validation error: No release defined in changelog")
		} else if (r.headers.release() || r.headers.change()) && !r.hasChangeDescriptions {
			err = errors.New("No change descriptions for " + r.headers.AsPath())
		}
	}
	return ast.WalkContinue, err
}

func (r *Renderer) visitHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var err error
	if entering {
		r.text.Reset()
	} else {
		var h Heading
		n := node.(*ast.Heading)
		switch n.Level {
		case 1:
			h, err = r.headers.ResetTo(titleHeading, r.text.String())
			if err != nil {
				break
			}

			title := h.(Title)
			err = r.validateDocumentHeading(title)
		case 2:
			if (r.headers.release() || r.headers.change()) && !r.hasChangeDescriptions {
				err = errors.New("No change descriptions for " + r.headers.AsPath())
			} else {
				h, err = r.headers.ResetTo(releaseHeading, r.text.String())
				if err != nil {
					break
				}

				release := h.(Release)
				err = r.validateReleaseHeading(release)
				if err != nil {
					break
				}

				if r.previousRelease.isRelease() && release.isRelease() {
					nextRelease := nextRelease(r.changes, release.version)
					if !r.previousRelease.version.EQ(nextRelease) {
						err = errors.New("Release '" + r.previousRelease.Name() + "' should have version " + nextRelease.String())
					}
				}

				if release.unreleased {
					r.h1Unreleased = true
				} else if !release.yanked {
					r.h1Released = true
				}
				r.hasChangeDescriptions = false
				r.changes = make(changeMap)
				r.previousRelease = release
			}
		case 3:
			if r.headers.title() {
				err = errors.New("Changes must be in a release " + r.headers.AsPath())
			} else if r.headers.change() && !r.hasChangeDescriptions {
				err = errors.New("No change descriptions for " + r.headers.AsPath())
			} else {
				h, err = r.headers.ResetTo(changeHeading, r.text.String())
				if err != nil {
					break
				}

				change := h.(Change)
				err = r.validateChangeHeading(change)
				if err != nil {
					break
				}
				r.hasChangeDescriptions = false
			}
		}
	}
	return ast.WalkContinue, err
}

func (r *Renderer) validateDocumentHeading(title Title) error { return nil }

func (r *Renderer) validateReleaseHeading(release Release) error {
	if release.unreleased {
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
		if release.yanked {
			if !r.h1Released && !r.h1Unreleased {
				return errors.New("Validation error: Changelog cannot start with a [YANKED] release, insert a release or a [Unreleased] first " + r.headers.AsPath())
			}
		}
		if (r.previousRelease.date != time.Time{}) {
			if r.previousRelease.date.Before(release.date) {
				return errors.New("Validation error: release '" + release.Name() + "' should be older than '" + r.previousRelease.Name() + "'")
			}
			if r.previousRelease.version.LTE(release.version) {
				return errors.New("Validation error: release '" + release.Name() + "' should be older than '" + r.previousRelease.Name() + "'")
			}
		}
	}
	return nil
}

func (r *Renderer) validateChangeHeading(change Change) error {
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
	if r.headers.change() {
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
