package query

import (
	"fmt"
	"strings"
)

type parserConfiguration struct {
	name     string
	elements expectedElements
}
type expectedElements map[string]expectedElement
type expectedElement struct {
	isScalar     bool
	enter, exit  project
	queryFactory queryFactory
}
type parsedElement struct {
	selector     string
	isRecursive  bool
	queryFactory queryFactory
}

type queryFactory func(selector string, isRecursive bool, queryElements []string) (Query, parsedElement, error)

func (expectedElements parserConfiguration) parseElement(queryElements []string) (element parsedElement, projection projections, err error) {
	name, selector, isScalar, isRecursive, err := parseName(queryElements[0])
	if err != nil {
		return parsedElement{}, projections{}, err
	}

	if expectedElement, ok := expectedElements.elements[name]; ok {
		if expectedElement.isScalar {
			if !isScalar {
				return parsedElement{}, projections{}, fmt.Errorf("%q is a scalar attribute", name)
			}
			if len(queryElements) != 1 {
				return parsedElement{}, projections{}, fmt.Errorf("no further query element allowed after %q", name)
			}
		} else if isScalar {
			return parsedElement{}, projections{}, fmt.Errorf("%q is a collection attribute", name)
		}
		return parsedElement{selector, isRecursive, expectedElement.queryFactory},
			projections{expectedElement.enter, expectedElement.exit, !isScalar && selector == ""},
			nil
	}

	return parsedElement{}, projections{}, fmt.Errorf("query attribute not recognized %q for a %q", name, expectedElements.name)
}

func parseName(element string) (name, selector string, isScalar, isRecursive bool, err error) {
	isRecursive = strings.HasSuffix(element, "/")
	openBracketIndex := strings.Index(element, "[")
	closeBracketIndex := strings.Index(element, "]")

	switch {
	case openBracketIndex != -1:
		if closeBracketIndex < openBracketIndex {
			err = fmt.Errorf("missing closing bracket in %q", element)
			return
		}

		name = element[:openBracketIndex]
		selector = element[openBracketIndex+1 : closeBracketIndex]
		isScalar = false
	case closeBracketIndex != -1:
		err = fmt.Errorf("missing opening bracket in %q", element)
		return
	default:
		if isRecursive {
			err = fmt.Errorf("recursion '/' not supported for scalar %q", element)
			return
		}

		name = element
		isScalar = true
	}
	return
}
