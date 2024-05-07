package query

import (
	"fmt"
	"strconv"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
)

func releaseQueryFactory(selector string, _ bool, queryElements []string) (Query, parsedElement, error) {
	i, err := strconv.Atoi(selector)
	if err != nil {
		return nil, parsedElement{}, fmt.Errorf("query release selector %q parsing error: %v", selector, err)
	}

	if len(queryElements) == 0 {
		return &releaseQuery{
			projections{func(of output.Format, h changelog.Heading) {
				if h, ok := h.(changelog.Release); ok {
					of.SetField("version", h.Version())
					of.SetField("date", h.Date())
				}
			}, nil, false,
			}, 0, i,
		}, parsedElement{}, nil
	}

	pe, projection, err := releaseParserConfiguration().parseElement(queryElements)
	if err != nil {
		return nil, parsedElement{}, err
	}

	return &releaseQuery{projection, 0, i}, pe, nil
}

func releaseParserConfiguration() parserConfiguration {
	return parserConfiguration{"release", expectedElements{
		"changes": {false, nil, nil, changeQueryFactory},
		"date": {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.Date())
			}
		}, nil, nil},
		"label": {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.Label())
			}
		}, nil, nil},
		"status": {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				switch {
				case !h.HasBeenReleased():
					of.Set("unreleased")
				case h.HasBeenYanked():
					of.Set("yanked")
				case h.IsPrerelease():
					of.Set("prereleased")
				default:
					of.Set("released")
				}
			}
		}, nil, nil},
		"title": {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.DisplayTitle())
			}
		}, nil, nil},
		"version": {true, func(of output.Format, h changelog.Heading) {
			if h, ok := h.(changelog.Release); ok {
				of.Set(h.Version())
			}
		}, nil, nil},
	}}
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
