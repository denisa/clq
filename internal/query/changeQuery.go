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
		_ = qe.newChangeItemQuery("", nil)
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Change); ok {
				rc.setField("title", h.Title())
				rc.array("changes")
			}
		}
		return nil
	}

	if strings.HasPrefix(queryElements[0], "descriptions[") && strings.HasSuffix(queryElements[0], "]") {
		if err := qe.newChangeItemQuery(selector(queryElements[0]), queryElements[1:]); err != nil {
			return err
		}
		return nil
	}

	switch queryElements[0] {
	default:
		return fmt.Errorf("Query attribute not recognized %q for a \"change\"", queryElements[0])
	case "title":
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Change); ok {
				rc.set(h.Title())
			}
		}
	}
	return nil
}

type changeQuery struct {
	projections
}

func (q *changeQuery) Enter(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.Change); !ok {
		return false, nil
	}
	return true, q.enter
}

func (q *changeQuery) Exit(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.Change); !ok {
		return false, nil
	}
	return true, q.exit
}
