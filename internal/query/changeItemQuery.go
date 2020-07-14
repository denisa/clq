package query

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newChangeItemQuery(name string, queryElements []string) error {
	if name != "" {
		return fmt.Errorf("Query change selector %q not yet supported", name)
	}
	if len(queryElements) > 0 {
		return fmt.Errorf("Query attribue selector %v not yet supported", queryElements)
	}
	queryMe := &changeItemQuery{enter: func(h changelog.ChangeItem) string {
		return fmt.Sprintf("%q,", h.Title())
	}}
	qe.queries = append(qe.queries, queryMe)
	return nil
}

type changeItemQuery struct {
	enter func(changelog.ChangeItem) string
	exit  func(changelog.ChangeItem) string
}

func (q *changeItemQuery) Enter(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.ChangeItem)
	if !ok {
		return false
	}

	if q.enter != nil {
		_, _ = w.WriteString(q.enter(h))
	}
	return true
}

func (q *changeItemQuery) Exit(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.ChangeItem)
	if !ok {
		return false
	}

	if q.exit != nil {
		_, _ = w.WriteString(q.exit(h))
	}
	return true
}
