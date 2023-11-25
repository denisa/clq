package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

func (qe *QueryEngine) newIntroductionQuery(queryElements []string) error {
	queryMe := &changelogQuery{}
	qe.queries = append(qe.queries, queryMe)

	elementName, elementSelector, elementIsList, elementIsRecursive, err := parseName(queryElements[0])
	if err != nil {
		return err
	}

	switch elementName {
	default:
		return fmt.Errorf("Query attribute not recognized %q for a \"introduction\"", elementName)
	case "releases":
		if err := elementIsCollection(elementName, elementIsList); err != nil {
			return err
		}
		if err := qe.newReleaseQuery(elementSelector, elementIsRecursive, queryElements[1:]); err != nil {
			return err
		}
	case "title":
		if err := elementIsFinal(elementName, elementIsList, queryElements[1:]); err != nil {
			return err
		}
		queryMe.enter = func(of output.OutputFormat, h changelog.Heading) {
			if h, ok := h.(changelog.Introduction); ok {
				of.Set(h.DisplayTitle())
			}
		}
	}
	return nil
}

type changelogQuery struct {
	projections
}

func (q *changelogQuery) isCollection() bool {
	return q.collection
}

func (q *changelogQuery) Accept(heading changelog.Heading) bool {
	_, ok := heading.(changelog.Introduction)
	return ok
}

func (q *changelogQuery) Enter(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.enter
}

func (q *changelogQuery) Exit(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.exit
}
