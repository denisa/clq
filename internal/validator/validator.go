package validator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/config"
	"github.com/denisa/clq/internal/semver"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// A Validator struct is an implementation of renderer.NodeRenderer that validates
// a changelog.
type Validator struct {
	release                  bool
	changeKind               *changelog.ChangeKind
	text                     strings.Builder
	hasIntroductionHeading   bool
	h1Released, h1Unreleased bool
	changes                  changelog.ChangeMap
	hasChangeDescriptions    bool
	changelog                *changelog.Changelog
	previousRelease          changelog.Release
}

func NewValidator(config config.Config) renderer.NodeRenderer {
	hf := changelog.NewHeadingFactory(config.ChangeKind())

	r := &Validator{
		release:    config.IsRelease(),
		changeKind: config.ChangeKind(),
		changes:    make(changelog.ChangeMap),
		changelog:  changelog.NewChangelog(hf),
	}

	if ok, listeners := config.Listeners(); ok {
		r.changelog.Listener(listeners)
	}

	return r
}

func (r *Validator) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDocument, r.visitDocument)
	reg.Register(ast.KindHeading, r.visitHeading)
	reg.Register(ast.KindList, r.visitList)
	reg.Register(ast.KindListItem, r.visitListItem)

	reg.Register(ast.KindAutoLink, r.visitAutoLink)
	reg.Register(ast.KindImage, r.visitImage)
	reg.Register(ast.KindLink, r.visitLink)
	reg.Register(ast.KindText, r.visitText)
}

func (r *Validator) visitDocument(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if !r.h1Released && !r.h1Unreleased {
			return ast.WalkStop, fmt.Errorf("validation error: No release defined in changelog")
		}
		if (r.changelog.Release() || r.changelog.Change()) && !r.hasChangeDescriptions {
			return ast.WalkStop, fmt.Errorf("no change descriptions for %v", r.changelog)
		}
		r.changelog.Close()
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitHeading(_ util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.text.Reset()
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Heading)
	if !r.hasIntroductionHeading && n.Level > 1 {
		return ast.WalkStop, fmt.Errorf("validation error: Introduction’s title must be defined")
	}

	switch n.Level {
	default:
		return ast.WalkStop, fmt.Errorf("validation error: Heading level %d not supported", n.Level)
	case 1:
		return r.visitHeading1()
	case 2:
		return r.visitHeading2()
	case 3:
		return r.visitHeading3()
	}
}

func (r *Validator) visitHeading1() (ast.WalkStatus, error) {
	_, err := r.changelog.Section(changelog.IntroductionHeading, r.text.String())
	if err != nil {
		return ast.WalkStop, err
	}
	r.hasIntroductionHeading = true
	return ast.WalkContinue, nil
}

func (r *Validator) visitHeading2() (ast.WalkStatus, error) {
	if (r.changelog.Release() || r.changelog.Change()) && !r.hasChangeDescriptions {
		return ast.WalkStop, fmt.Errorf("no change descriptions for %v", r.changelog)
	}
	h, err := r.changelog.Section(changelog.ReleaseHeading, r.text.String())
	if err != nil {
		return ast.WalkStop, err
	}

	release := h.(changelog.Release)
	if err := r.validateReleaseHeading(release); err != nil {
		return ast.WalkStop, err
	}

	if r.previousRelease.IsRelease() && release.IsRelease() {
		changeKind := r.changeKind
		increment, trigger := changeKind.IncrementFor(r.changes)
		if increment == semver.Build {
			return ast.WalkStop, fmt.Errorf("release %q cannot have only build-level changes because it is not the initial release", r.previousRelease.Title())
		}
		if err := r.previousRelease.IsNewerThan(release); err != nil {
			return ast.WalkStop, err
		}
		nextRelease := release.NextRelease(increment)
		if !r.previousRelease.ReleaseIs(nextRelease) {
			if release.IsMajorVersionZero() && increment == semver.Major {
				nextMinorRelease := release.NextRelease(semver.Minor)
				if !r.previousRelease.ReleaseIs(nextMinorRelease) {
					return ast.WalkStop, fmt.Errorf("release %q should have version %v or %v because of %q", r.previousRelease.Title(), nextMinorRelease, nextRelease, trigger)
				}
			} else {
				return ast.WalkStop, fmt.Errorf("release %q should have version %v because of %q", r.previousRelease.Title(), nextRelease, trigger)
			}
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
	return ast.WalkContinue, nil
}

func (r *Validator) validateReleaseHeading(release changelog.Release) error {
	if release.HasBeenReleased() {
		if release.HasBeenYanked() {
			if !r.h1Released && !r.h1Unreleased {
				return fmt.Errorf("validation error: Changelog cannot start with a \"[YANKED]\" release, insert a release or a \"[Unreleased]\" first %v", r.changelog)
			}
		}
	} else {
		if r.release {
			return fmt.Errorf("validation error: \"[Unreleased]\" not supported in release mode %v", r.changelog)
		}
		if r.h1Unreleased {
			return fmt.Errorf("validation error: Multiple \"[Unreleased]\" not supported %v", r.changelog)
		}
		if r.h1Released {
			return fmt.Errorf("validation error: \"[Unreleased]\" must come before any release %v", r.changelog)
		}
	}
	return nil
}

func (r *Validator) visitHeading3() (ast.WalkStatus, error) {
	if r.changelog.Introduction() {
		return ast.WalkStop, fmt.Errorf("changes must be in a release %v", r.changelog)
	}
	if r.changelog.Change() && !r.hasChangeDescriptions {
		return ast.WalkStop, fmt.Errorf("no change descriptions for %v", r.changelog)
	}

	h, err := r.changelog.Section(changelog.ChangeHeading, r.text.String())
	if err != nil {
		return ast.WalkStop, err
	}

	change := h.(changelog.Change)
	if err := r.validateChangeHeading(change); err != nil {
		return ast.WalkStop, err
	}
	r.hasChangeDescriptions = false
	return ast.WalkContinue, nil
}

func (r *Validator) validateChangeHeading(change changelog.Change) error {
	if r.changes[change.Title()] {
		return fmt.Errorf("validation error: Multiple headings %q not supported %v", change.Title(), r.changelog)
	}
	r.changes[change.Title()] = true
	return nil
}

func (r *Validator) visitAutoLink(_ util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if entering {
		r.text.WriteString("<")
		url := n.URL(source)
		if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
			r.text.WriteString("mailto:")
		}
		r.text.Write(url)
		r.text.WriteString(">")
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitLink(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		r.text.WriteString("[")
	} else {
		r.text.WriteString("](")
		r.text.Write(n.Destination)
		_ = w.WriteByte('"')
		if n.Title != nil {
			r.text.WriteString(" \"")
			r.text.Write(n.Title)
			r.text.WriteString("\"")
		}
		r.text.WriteString(")")
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitList(_ util.BufWriter, _ []byte, _ ast.Node, _ bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Validator) visitListItem(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.text.Reset()
		return ast.WalkContinue, nil
	}
	if r.changelog.Change() {
		_, err := r.changelog.Section(changelog.ChangeDescription, r.text.String())
		if err != nil {
			return ast.WalkStop, err
		}
		r.hasChangeDescriptions = true
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitImage(_ util.BufWriter, _ []byte, _ ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Validator) visitText(_ util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Text)
	if entering {
		segment := n.Segment
		value := segment.Value(source)
		r.text.Write(value)
		if !n.IsRaw() {
			if n.HardLineBreak() {
				r.text.WriteString("  \n")
			} else if n.SoftLineBreak() {
				r.text.WriteString(" ")
			}
		}
	}
	return ast.WalkContinue, nil
}
