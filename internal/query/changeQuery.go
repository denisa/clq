package query

import (
	"fmt"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newChangeQuery(name string, queryElements []string) error {
	if name != "" {
		return fmt.Errorf("Query change selector %q not yet supported", name)
	}

	queryMe := &changeQuery{}
	qe.queries = append(qe.queries, queryMe)

	if len(queryElements) == 0 {
		if err := qe.newChangeItemQuery("", nil); err != nil {
			return err
		}
		queryMe.enter = func(c changelog.Change) string {
			return fmt.Sprintf("{\"kind\":%q, \"changes\":[", c.Title())
		}
		queryMe.exit = func(c changelog.Change) string {
			return "]}"
		}
		return nil
	}

	return fmt.Errorf("Query attribute not recognized %q", name)
}

type changeQuery struct {
	enter func(changelog.Change) string
	exit  func(changelog.Change) string
}

func (q *changeQuery) Enter(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Change)
	if !ok {
		return false
	}

	if q.enter != nil {
		_, _ = w.WriteString(q.enter(h))
	}
	return true
}

func (q *changeQuery) Exit(w *strings.Builder, heading changelog.Heading) bool {
	h, ok := heading.(changelog.Change)
	if !ok {
		return false
	}

	if q.exit != nil {
		_, _ = w.WriteString(q.exit(h))
	}
	return true
}
