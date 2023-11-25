package config

import (
	"github.com/denisa/clq/internal/changelog"
)

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

// withChangeKind is a functional option that allow you to set the supported
// change kinds.
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
// Listener.
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

// WithRelease is a functional option that allow you to set the Release mode to
// the Validator.
func WithRelease(release bool) interface {
	Option
} {
	return &withRelease{release}
}
