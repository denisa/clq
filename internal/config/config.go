// Package config provides methods to track choices made on the command line
package config

import (
	"github.com/denisa/clq/internal/changelog"
)

// A Config struct has configurations for the Validator.
type Config struct {
	release    bool
	listener   changelog.Listener
	changeKind *changelog.ChangeKind
}

// NewConfig builds a new Config with all Options.
func NewConfig(opts ...Option) Config {
	config := Config{}
	for _, opt := range opts {
		opt.SetValidationOption(&config)
	}
	return config
}

func (c Config) IsRelease() bool {
	return c.release
}

func (c Config) Listeners() (bool, changelog.Listener) {
	return c.listener != nil, c.listener
}

func (c Config) ChangeKind() *changelog.ChangeKind {
	return c.changeKind
}
