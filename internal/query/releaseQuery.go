package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newReleaseQuery(name string, queryElements []string) error {
	i, err := strconv.Atoi(name)
	if err != nil {
		return fmt.Errorf("Query release selector %q parsing error: %v", name, err)
	}

	queryMe := &releaseQuery{index: i}
	qe.queries = append(qe.queries, queryMe)

	if len(queryElements) == 0 {
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				rc.setField("version", h.Version())
				rc.setField("date", h.Date())
			}
		}
		return nil
	}

	if strings.HasPrefix(queryElements[0], "changes[") && strings.HasSuffix(queryElements[0], "]") {
		if err := qe.newChangeQuery(selector(queryElements[0]), queryElements[1:]); err != nil {
			return err
		}
		return nil
	}

	switch queryElements[0] {
	default:
		return fmt.Errorf("Query attribute not recognized %q for a \"release\"", queryElements[0])
	case "date":
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				rc.set(h.Date())
			}
		}
	case "label":
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				rc.set(h.Label())
			}
		}
	case "status":
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				if !h.HasBeenReleased() {
					rc.set("unreleased")
				} else if h.HasBeenYanked() {
					rc.set("yanked")
				} else if h.IsPrerelease() {
					rc.set("prereleased")
				} else {
					rc.set("released")
				}
			}
		}
	case "version":
		queryMe.enter = func(rc ResultCollector, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				rc.set(h.Version())
			}
		}
	}

	return nil
}

type releaseQuery struct {
	projections
	cursor int
	index  int
}

func (q *releaseQuery) Enter(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.Release); !ok {
		return false, nil
	}
	selected := q.cursor == q.index
	q.cursor++

	if !selected {
		return false, nil
	}
	return true, q.enter
}

func (q *releaseQuery) Exit(heading changelog.Heading) (bool, Project) {
	if _, ok := heading.(changelog.Release); !ok {
		return false, nil
	}
	return true, q.exit
}
