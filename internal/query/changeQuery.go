package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

const (
	jsonNameDescriptions string = "descriptions"
	jsonNameTitle        string = "title"
)

func changeQueryFactory(selector string, isRecursive bool, queryElements []string) (Query, parsedElement, error) {
	if selector != "" {
		return nil, parsedElement{}, fmt.Errorf("query change selector %q not yet supported", selector)
	}

	queryMe := &changeQuery{}
	queryMe.collection = true

	parsedElement := parsedElement{}

	if len(queryElements) == 0 {
		if isRecursive {
			queryMe.enter = func(of output.Format, h changelog.Heading) {
				if h, ok := h.(changelog.Change); ok {
					of.SetField(jsonNameTitle, h.DisplayTitle())
					of.Array(jsonNameDescriptions)
				}
			}
			parsedElement.queryFactory = changeItemQueryFactory
			parsedElement.isRecursive = true
		} else {
			queryMe.enter = func(of output.Format, h changelog.Heading) {
				if h, ok := h.(changelog.Change); ok {
					of.SetField(jsonNameTitle, h.DisplayTitle())
				}
			}
		}
		return queryMe, parsedElement, nil
	}

	parsedElement, projection, err := changeParserConfiguration().parseElement(queryElements)
	if err != nil {
		return nil, parsedElement, err
	}

	return &changeQuery{projection}, parsedElement, nil
}

func changeParserConfiguration() parserConfiguration {
	return parserConfiguration{"change", expectedElements{
		jsonNameDescriptions: {false, nil, nil, changeItemQueryFactory},
		jsonNameTitle: {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Change); ok {
				of.Set(h.DisplayTitle())
			}
		}, nil, nil}}}
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
