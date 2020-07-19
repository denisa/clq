package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
)

func (qe *QueryEngine) newChangeQuery(selector string, isRecursive bool, queryElements []string) error {
	if selector != "" {
		return fmt.Errorf("Query change selector %q not yet supported", selector)
	}

	queryMe := &changeQuery{}
	queryMe.collection = true
	qe.queries = append(qe.queries, queryMe)

	if len(queryElements) == 0 {
		if isRecursive {
			_ = qe.newChangeItemQuery("", true, nil)
			queryMe.enter = func(rc resultCollector, h changelog.Heading) {
				if h, ok := h.(changelog.Change); ok {
					rc.setField("title", h.Title())
					rc.array("descriptions")
				}
			}
		} else {
			queryMe.enter = func(rc resultCollector, h changelog.Heading) {
				if h, ok := h.(changelog.Change); ok {
					rc.setField("title", h.Title())
				}
			}
		}
		return nil
	}

	elementName, elementSelector, elementIsList, elementIsRecursive, err := parseName(queryElements[0])
	if err != nil {
		return err
	}

	switch elementName {
	default:
		return fmt.Errorf("Query attribute not recognized %q for a \"change\"", elementName)
	case "descriptions":
		if err := elementIsCollection(elementName, elementIsList); err != nil {
			return err
		}
		if err := qe.newChangeItemQuery(elementSelector, elementIsRecursive, queryElements[1:]); err != nil {
			return err
		}
	case "title":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(rc resultCollector, h changelog.Heading) {
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

func (q *changeQuery) isCollection() bool {
	return q.collection
}

func (q *changeQuery) Accept(heading changelog.Heading) bool {
	_, ok := heading.(changelog.Change)
	return ok
}

func (q *changeQuery) Enter(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.enter
}

func (q *changeQuery) Exit(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.exit
}
