package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
)

type OutputFormat interface {
	Result() string
	Open(heading changelog.Heading)
	Close(heading changelog.Heading)
}

type ResultCollector interface {
	set(value string)
	setField(name string, value string)
	array(name string)
}

type Project func(rc ResultCollector, heading changelog.Heading)

type projections struct {
	enter Project
	exit  Project
}

func newOutputFormat(formatName string) (OutputFormat, error) {
	switch formatName {
	case "json":
		return &jsonResultCollector{}, nil
	case "md":
		return &mdResultCollector{}, nil
	default:
		return nil, fmt.Errorf("Unrecognized output format %q. Supported format: \"json\", \"md\"", formatName)
	}
}
