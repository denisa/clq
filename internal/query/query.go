// query is a simple streaming query language to identify parts of a CHANGELOG.
// The query language is inspired by xpath: it consists of a sequence of individual
// queries, one for each Heading traversed from the root, the changelog itself, to the
// desired element.
package query

import (
	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

// Query is the query for a single heading.
type Query interface {
	// enter returns true if this query element accepts this kind of heading.
	Accept(heading changelog.Heading) bool
	// enter returns true if the given heading fulfils the query expression.
	// enter might write to the buffer part of the query result.
	Enter(heading changelog.Heading) (bool, project)
	// exit returns true if the given heading fulfils the query expression.
	// exit might write to the buffer part of the query result.
	Exit(heading changelog.Heading) (bool, project)
	// isCollection returns true is this query produces a collection of results
	isCollection() bool
}

// a project function projects the desired part of the heading in the output.Format.
type project func(of output.Format, heading changelog.Heading)

// projections is a base type for all queries
type projections struct {
	enter, exit project
	collection  bool
}
