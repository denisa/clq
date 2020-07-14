// query is a simple streaming query language to identify parts of a CHANGELOG.
// The query language is inspired by xpath: it consists of a sequence of individual
// queries, one for each Heading traversed from the root, the changelog itself, to the
// desired element.
package query

import (
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

// Query is the query for a single heading.
type Query interface {
	// enter returns true if the given heading fulfils the query expression.
	// enter might write to the buffer part of the query result.
	Enter(w *strings.Builder, heading changelog.Heading) bool
	// exit returns true if the given heading fulfils the query expression.
	// exit might write to the buffer part of the query result.
	Exit(w *strings.Builder, heading changelog.Heading) bool
}
