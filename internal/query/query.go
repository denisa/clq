// query is a simple streaming query language to identify parts of a CHANGELOG.
// The query language is inspired by xpath: it consists of a sequence of individual
// queries, one for each Heading traversed from the root, the changelog itself, to the
// desired element.
package query

import (
	"github.com/denisa/clq/internal/changelog"
)

// Query is the query for a single heading.
type Query interface {
	// enter returns true if the given heading fulfils the query expression.
	// enter might write to the buffer part of the query result.
	Enter(heading changelog.Heading) (bool, Project)
	// exit returns true if the given heading fulfils the query expression.
	// exit might write to the buffer part of the query result.
	Exit(heading changelog.Heading) (bool, Project)
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
