package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newChangeItemQuery(name string, queryElements []string) error {
	if name != "" {
		return fmt.Errorf("Query change selector %q not yet supported", name)
	}
	if len(queryElements) > 0 {
		return fmt.Errorf("Query attribute selector %v not yet supported", queryElements)
	}
	queryMe := &changeItemQuery{}
	queryMe.exit = func(rc ResultCollector, h changelog.Heading) {
		if h, ok := h.(changelog.ChangeItem); ok {
			rc.set(h.Title())
		}
	}
	qe.queries = append(qe.queries, queryMe)
	return nil
}

type changeItemQuery struct {
	projections
}

func (q *changeItemQuery) Enter(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.ChangeItem); !ok {
		return false, nil
	}
	return true, q.enter
}

func (q *changeItemQuery) Exit(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.ChangeItem); !ok {
		return false, nil
	}
	return true, q.exit
}
