// query is a simple streaming query language to identify parts of a CHANGELOG.
// The query language is inspaired by xpath: it consists of a sequence of individual
// queries, one for each Heading traversed from the root, the changelog itself, to the
// desired element.
package query

import (
	"github.com/denisa/clq/internal/changelog"
	"github.com/yuin/goldmark/util"
)

// Query is the query for a single heading.
type Query interface {
	// select returns true if the given heading fulfils the query expression.
	// select might write to the buffer part of the query result.
	Select(w util.BufWriter, heading changelog.Heading) bool
}
