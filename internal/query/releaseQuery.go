package query

import (
	"fmt"
	"strconv"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

func (qe *QueryEngine) newReleaseQuery(selector string, _ bool, queryElements []string) error {
	i, err := strconv.Atoi(selector)
	if err != nil {
		return fmt.Errorf("Query release selector %q parsing error: %v", selector, err)
	}

	queryMe := &releaseQuery{index: i}
	qe.queries = append(qe.queries, queryMe)

	if len(queryElements) == 0 {
		queryMe.enter = func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.SetField("version", h.Version())
				of.SetField("date", h.Date())
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
		return fmt.Errorf("Query attribute not recognized %q for a \"release\"", elementName)
	case "changes":
		if err := elementIsCollection(elementName, elementIsList); err != nil {
			return err
		}
		if err := qe.newChangeQuery(elementSelector, elementIsRecursive, queryElements[1:]); err != nil {
			return err
		}
	case "date":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.Date())
			}
		}
	case "label":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.Label())
			}
		}
	case "status":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				if !h.HasBeenReleased() {
					of.Set("unreleased")
				} else if h.HasBeenYanked() {
					of.Set("yanked")
				} else if h.IsPrerelease() {
					of.Set("prereleased")
				} else {
					of.Set("released")
				}
			}
		}
	case "title":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.DisplayTitle())
			}
		}
	case "version":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.Version())
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

func (q *releaseQuery) isCollection() bool {
	return q.collection
}

func (q *releaseQuery) Accept(heading changelog.Heading) bool {
	_, ok := heading.(changelog.Release)
	return ok
}

func (q *releaseQuery) Enter(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	selected := q.cursor == q.index
	q.cursor++

	if !selected {
		return false, nil
	}
	return true, q.enter
}

func (q *releaseQuery) Exit(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.exit
}
