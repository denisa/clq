package validator

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

const semver string = `(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`

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
	reg.Register(ast.KindParagraph, r.visitParagraph)
	reg.Register(ast.KindTextBlock, r.visitTextBlock)
	reg.Register(ast.KindThematicBreak, r.visitThematicBreak)

	reg.Register(ast.KindAutoLink, r.visitAutoLink)
	reg.Register(ast.KindCodeSpan, r.visitCodeSpan)
	reg.Register(ast.KindEmphasis, r.visitEmphasis)
	reg.Register(ast.KindImage, r.visitImage)
	reg.Register(ast.KindLink, r.visitLink)
	reg.Register(ast.KindRawHTML, r.visitRawHTML)
	reg.Register(ast.KindText, r.visitText)
	reg.Register(ast.KindString, r.visitString)
}

func (r *Renderer) visitDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var err error
	if !entering {
		if !r.h1Released && !r.h1Unreleased {
			err = errors.New("Validation error: No release defined in changelog")
		} else if (r.headers.release() || r.headers.change()) && !r.hasChangeDescriptions {
			err = errors.New("No change descriptions for " + r.headers.asPath())
		}
	}
	return ast.WalkContinue, err
}

func (r *Renderer) visitHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var err error
	if entering {
		r.text.Reset()
	} else {
		n := node.(*ast.Heading)
		switch n.Level {
		case 1:
			r.headers.resetTo(0, r.text.String())
			err = r.validateDocumentHeading()
		case 2:
			if (r.headers.release() || r.headers.change()) && !r.hasChangeDescriptions {
				err = errors.New("No change descriptions for " + r.headers.asPath())
			} else {
				r.headers.resetTo(1, r.text.String())
				err = r.validateReleaseHeading()
				r.hasChangeDescriptions = false
			}
		case 3:
			if r.headers.title() {
				err = errors.New("Changes must be in a release " + r.headers.asPath())
			} else if r.headers.change() && !r.hasChangeDescriptions {
				err = errors.New("No change descriptions for " + r.headers.asPath())
			} else {
				r.headers.resetTo(2, r.text.String())
				err = r.validateChangeHeading()
				r.hasChangeDescriptions = false
			}
		}
	}
	return ast.WalkContinue, err
}

func (r *Renderer) validateDocumentHeading() error { return nil }

var yankedRelease = regexp.MustCompile(`^` + semver + `\s+-\s+(?P<date>\d\d\d\d-\d\d-\d\d)\s+\[\s*YANKED\s*]?$`)
var release = regexp.MustCompile(`^\[\s*` + semver + `\s*\]\s+-\s+(?P<date>\d\d\d\d-\d\d-\d\d)(?:\s+(?P<label>.+))?$`)

func (r *Renderer) validateReleaseHeading() error {
	if matched, _ := regexp.MatchString(`^\[\s*Unreleased\s*\]$`, r.text.String()); matched {
		if r.Release {
			return errors.New("Validation error: [Unreleased] not supported in release mode " + r.headers.asPath())
		}
		if r.h1Unreleased {
			return errors.New("Validation error: Multiple [Unreleased] not supported " + r.headers.asPath())
		}
		if r.h1Released {
			return errors.New("Validation error: [Unreleased] must come before any release " + r.headers.asPath())
		}
		r.h1Unreleased = true
	} else if matches := yankedRelease.FindStringSubmatch(r.text.String()); matches != nil {
		if !r.h1Released && !r.h1Unreleased {
			return errors.New("Validation error: Changelog cannot start with a [YANKED] release, insert a release or a [Unreleased] first " + r.headers.asPath())
		}
		if _, err := time.Parse("2006-01-02", subexp(yankedRelease, matches, "date")); err != nil {
			return errors.New("Validation error: Illegal date (" + err.Error() + ")" + r.headers.asPath())
		}
	} else if matches := release.FindStringSubmatch(r.text.String()); matches != nil {
		if matched, _ := regexp.MatchString(`[\s*YANKED\s*]`, subexp(release, matches, "label")); matched {
			return errors.New("Validation error: the version of a [YANKED] release cannot stand between [...] " + r.headers.asPath())
		} else if _, err := time.Parse("2006-01-02", subexp(release, matches, "date")); err != nil {
			return errors.New("Validation error: Illegal date (" + err.Error() + ")" + r.headers.asPath())
		}
		r.h1Released = true
	}
	r.changes = make(changeMap)
	return nil
}

func subexp(exp *regexp.Regexp, matches []string, subexp string) string {
	for index, name := range exp.SubexpNames() {
		if index >= len(matches) {
			continue
		}

		if name == subexp {
			return matches[index]
		}
	}

	return ""
}

var changes = []string{"Added", "Removed", "Changed", "Deprecated", "Fixed", "Security"}

func (r *Renderer) validateChangeHeading() error {
	change := r.text.String()
	for _, val := range changes {
		if matched, _ := regexp.MatchString(`^`+val+`$`, change); matched {
			if r.changes[val] {
				return errors.New("Validation error: Multiple headings " + val + " not supported " + r.headers.asPath())
			}
			r.changes[val] = true
			return nil
		}
	}
	return errors.New("Validation error: Unknown change headings " + change + " not supported " + r.headers.asPath())
}

func (r *Renderer) visitBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.writeLines(w, source, node)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.writeLines(w, source, node)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.writeLines(w, source, node)
	}
	return ast.WalkContinue, nil
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

func (r *Renderer) visitParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitTextBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if _, ok := node.NextSibling().(ast.Node); ok && node.FirstChild() != nil {
			r.text.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) visitImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) visitRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func (r *Renderer) visitString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.String)
		r.text.Write(n.Value)
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, node ast.Node) {
	l := node.Lines().Len()
	for i := 0; i < l; i++ {
		line := node.Lines().At(i)
		r.text.Write(line.Value(source))
	}
}
