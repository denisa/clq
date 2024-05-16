package query

import (
	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

func introductionQueryFactory(_ string, _ bool, queryElements []string) (Query, parsedElement, error) {
	pe, projection, err := changelogParserConfiguration().parseElement(queryElements)
	if err != nil {
		return nil, parsedElement{}, err
	}

	return &changelogQuery{projection}, pe, nil
}

func changelogParserConfiguration() parserConfiguration {
	return parserConfiguration{"introduction", expectedElements{
		"releases": {false, nil, nil, releaseQueryFactory},
		"title": {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Introduction); ok {
				of.Set(h.DisplayTitle())
			}
		}, nil, nil},
	}}
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
