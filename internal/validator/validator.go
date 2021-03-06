package validator

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A Config struct has configurations for the Validator.
type Config struct {
	release    bool
	listener   changelog.Listener
	changeKind *changelog.ChangeKind
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{}
}

// An Option interface sets options for the Validator.
type Option interface {
	SetValidationOption(*Config)
}

// ------------- ChangeKind -------------
type withChangeKind struct {
	value *changelog.ChangeKind
}

func (o *withChangeKind) SetValidationOption(c *Config) {
	c.changeKind = o.value
}

// withChangeKind is a functional option that allow you to set the changelog event
// listener.
func WithChangeKind(changeKind *changelog.ChangeKind) interface {
	Option
} {
	return &withChangeKind{value: changeKind}
}

// ------------- Listener -------------
type withListener struct {
	value changelog.Listener
}

func (o *withListener) SetValidationOption(c *Config) {
	c.listener = o.value
}

// withListener is a functional option that allow you to set the changelog event
// listener.
func WithListener(listener changelog.Listener) interface {
	Option
} {
	return &withListener{value: listener}
}

// ------------- Release -------------
type withRelease struct {
	value bool
}

func (o *withRelease) SetValidationOption(c *Config) {
	c.release = o.value
}

// WithRelease is a functional option that allow you to set the release mode to
// the Validator.
func WithRelease(release bool) interface {
	Option
} {
	return &withRelease{release}
}

// A Validator struct is an implementation of renderer.NodeRenderer that validates
// a changelog.
type Validator struct {
	Config
	text                     strings.Builder
	h1Released, h1Unreleased bool
	changes                  changelog.ChangeMap
	hasChangeDescriptions    bool
	headers                  changelog.Changelog
	previousRelease          changelog.Release
}

func NewValidator(opts ...Option) renderer.NodeRenderer {
	r := &Validator{
		Config:  NewConfig(),
		changes: make(changelog.ChangeMap),
		headers: changelog.NewChangelog(),
	}

	for _, opt := range opts {
		opt.SetValidationOption(&r.Config)
	}

	if r.listener != nil {
		r.headers.Listener(r.listener)
	}

	return r
}

func (r *Validator) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDocument, r.visitDocument)
	reg.Register(ast.KindHeading, r.visitHeading)
	reg.Register(ast.KindList, r.visitList)
	reg.Register(ast.KindListItem, r.visitListItem)

	reg.Register(ast.KindImage, r.visitImage)
	reg.Register(ast.KindText, r.visitText)
}

func (r *Validator) visitDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if !r.h1Released && !r.h1Unreleased {
			return ast.WalkStop, fmt.Errorf("Validation error: No release defined in changelog")
		}
		if (r.headers.Release() || r.headers.Change()) && !r.hasChangeDescriptions {
			return ast.WalkStop, fmt.Errorf("No change descriptions for %v", r.headers)
		}
		r.headers.Close()
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.text.Reset()
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Heading)
	switch n.Level {
	case 1:
		_, err := r.headers.Section(changelog.IntroductionHeading, r.text.String())
		if err != nil {
			return ast.WalkStop, err
		}
		// no validation rules defined for the title...
	case 2:
		if (r.headers.Release() || r.headers.Change()) && !r.hasChangeDescriptions {
			if err := fmt.Errorf("No change descriptions for %v", r.headers); err != nil {
				return ast.WalkStop, err
			}
		} else {
			h, err := r.headers.Section(changelog.ReleaseHeading, r.text.String())
			if err != nil {
				return ast.WalkStop, err
			}

			release := h.(changelog.Release)
			if err := r.validateReleaseHeading(release); err != nil {
				return ast.WalkStop, err
			}

			if r.previousRelease.IsRelease() && release.IsRelease() {
				nextRelease := release.NextRelease(r.changeKind.IncrementFor(r.changes))
				if !r.previousRelease.ReleaseIs(nextRelease) {
					return ast.WalkStop, fmt.Errorf("Release %q should have version %v", r.previousRelease.Title(), nextRelease)
				}
			}

			if !release.HasBeenReleased() {
				r.h1Unreleased = true
			} else if !release.HasBeenYanked() {
				r.h1Released = true
			}
			r.hasChangeDescriptions = false
			r.changes = make(changelog.ChangeMap)
			r.previousRelease = release
		}
	case 3:
		if r.headers.Introduction() {
			return ast.WalkStop, fmt.Errorf("Changes must be in a release %v", r.headers)
		}
		if r.headers.Change() && !r.hasChangeDescriptions {
			return ast.WalkStop, fmt.Errorf("No change descriptions for %v", r.headers)
		}

		h, err := r.headers.Section(changelog.ChangeHeading, r.text.String())
		if err != nil {
			return ast.WalkStop, err
		}

		change := h.(changelog.Change)
		if err := r.validateChangeHeading(change); err != nil {
			return ast.WalkStop, err
		}
		r.hasChangeDescriptions = false
	}
	return ast.WalkContinue, nil
}

func (r *Validator) validateReleaseHeading(release changelog.Release) error {
	if !release.HasBeenReleased() {
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
		if release.HasBeenYanked() {
			if !r.h1Released && !r.h1Unreleased {
				return fmt.Errorf("Validation error: Changelog cannot start with a \"[YANKED]\" release, insert a release or a \"[Unreleased]\" first %v", r.headers)
			}
		}
		if r.previousRelease.HasBeenReleased() {
			if err := r.previousRelease.IsNewerThan(release); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Validator) validateChangeHeading(change changelog.Change) error {
	if err := r.changeKind.IsSupported(change.Title()); err != nil {
		return err
	}
	if r.changes[change.Title()] {
		return fmt.Errorf("Validation error: Multiple headings %q not supported %v", change.Title(), r.headers)
	}
	r.changes[change.Title()] = true
	return nil
}

func (r *Validator) visitList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Validator) visitListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.text.Reset()
		return ast.WalkContinue, nil
	}
	if r.headers.Change() {
		_, err := r.headers.Section(changelog.ChangeDescription, r.text.String())
		if err != nil {
			return ast.WalkStop, err
		}
		r.hasChangeDescriptions = true
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Text)
		segment := n.Segment
		r.text.Write(segment.Value(source))
	}
	return ast.WalkContinue, nil
}
