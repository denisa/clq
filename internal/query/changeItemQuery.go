package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

func changeItemQueryFactory(selector string, _ bool, queryElements []string) (Query, parsedElement, error) {
	if selector != "" {
		return nil, parsedElement{}, fmt.Errorf("query change selector %q not yet supported", selector)
	}
	if len(queryElements) > 0 {
		return nil, parsedElement{}, fmt.Errorf("query attribute selector %v not yet supported", queryElements)
	}
	queryMe := &changeItemQuery{}
	queryMe.collection = true
	queryMe.exit = func(of output.Format, h changelog.Heading) {
		if h, ok := h.(changelog.ChangeItem); ok {
			of.Set(h.DisplayTitle())
		}
	}
	return queryMe, parsedElement{}, nil
}

type changeItemQuery struct {
	projections
}

func (q *changeItemQuery) isCollection() bool {
	return q.collection
}

func (q *changeItemQuery) Accept(heading changelog.Heading) bool {
	_, ok := heading.(changelog.ChangeItem)
	return ok
}

func (q *changeItemQuery) Enter(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.enter
}

func (q *changeItemQuery) Exit(heading changelog.Heading) (bool, project) {
	if !q.Accept(heading) {
		return false, nil
	}
	return true, q.exit
}
