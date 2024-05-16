package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseElementUnkownAttributeError(t *testing.T) {
	_, _, errParseElement := parserConfiguration{"failing", expectedElements{
		"supported": {true, nil, nil, nil},
	}}.parseElement([]string{"unsupported"})
	require.EqualError(t, errParseElement, "query attribute not recognized \"unsupported\" for a \"failing\"")
}

func TestParseElementAttributeShouldBeScalarError(t *testing.T) {
	_, _, errParseElement := parserConfiguration{"failing", expectedElements{
		"supported": {true, nil, nil, nil},
	}}.parseElement([]string{"supported[]"})
	require.EqualError(t, errParseElement, "\"supported\" is a scalar attribute")
}

func TestParseElementScalarAttributeShouldEndQueyError(t *testing.T) {
	_, _, errParseElement := parserConfiguration{"failing", expectedElements{
		"supported": {true, nil, nil, nil},
	}}.parseElement([]string{"supported", "unsupported"})
	require.EqualError(t, errParseElement, "no further query element allowed after \"supported\"")
}

func TestParseElementAttributeShouldNotBeScalarError(t *testing.T) {
	_, _, errParseElement := parserConfiguration{"failing", expectedElements{
		"supported": {false, nil, nil, nil},
	}}.parseElement([]string{"supported"})
	require.EqualError(t, errParseElement, "\"supported\" is a collection attribute")
}

func TestParseElementScalarAttribute(t *testing.T) {
	parsedElement, projection, errParseElement := parserConfiguration{"succeed", expectedElements{
		"scalar":     {true, nil, nil, nil},
		"collection": {false, nil, nil, nil},
	}}.parseElement([]string{"scalar"})

	assertions := require.New(t)
	assertions.NoError(errParseElement)
	assertions.False(parsedElement.isRecursive)
	assertions.Empty(parsedElement.selector)
	assertions.False(projection.collection)
}

func TestParseElementCollectionAttributSelectorNotRecursive(t *testing.T) {
	parsedElement, projection, errParseElement := parserConfiguration{"succeed", expectedElements{
		"scalar":     {true, nil, nil, nil},
		"collection": {false, nil, nil, nil},
	}}.parseElement([]string{"collection[three]"})

	assertions := require.New(t)
	assertions.NoError(errParseElement)
	assertions.False(parsedElement.isRecursive)
	assertions.Equal(parsedElement.selector, "three")
	assertions.False(projection.collection)
}

func TestParseElementCollectionAttributeSelectorNotRecursive(t *testing.T) {
	parsedElement, projection, errParseElement := parserConfiguration{"succeed", expectedElements{
		"scalar":     {true, nil, nil, nil},
		"collection": {false, nil, nil, nil},
	}}.parseElement([]string{"collection[three]/"})

	assertions := require.New(t)
	assertions.NoError(errParseElement)
	assertions.True(parsedElement.isRecursive)
	assertions.Equal(parsedElement.selector, "three")
	assertions.False(projection.collection)
}

func TestParseElementCollectionAttributeNoSelectorNotRecursive(t *testing.T) {
	parsedElement, projection, errParseElement := parserConfiguration{"succeed", expectedElements{
		"scalar":     {true, nil, nil, nil},
		"collection": {false, nil, nil, nil},
	}}.parseElement([]string{"collection[]"})

	assertions := require.New(t)
	assertions.NoError(errParseElement)
	assertions.False(parsedElement.isRecursive)
	assertions.Empty(parsedElement.selector)
	assertions.True(projection.collection)
}

func TestParseElementCollectionAttributeNoSelectorRecursive(t *testing.T) {
	parsedElement, projection, errParseElement := parserConfiguration{"succeed", expectedElements{
		"scalar":     {true, nil, nil, nil},
		"collection": {false, nil, nil, nil},
	}}.parseElement([]string{"collection[]/"})

	assertions := require.New(t)
	assertions.NoError(errParseElement)
	assertions.True(parsedElement.isRecursive)
	assertions.Empty(parsedElement.selector)
	assertions.True(projection.collection)
}

func TestParseNameFormatError(t *testing.T) {
	testcases := []struct{ element, error string }{
		{"changes[", "missing closing bracket in \"changes[\""},
		{"changes]", "missing opening bracket in \"changes]\""},
		{"changes][", "missing closing bracket in \"changes][\""},
		{"changes/", "recursion '/' not supported for scalar \"changes/\""},
	}
	for _, testcase := range testcases {
		t.Run(testcase.element, func(t *testing.T) {
			_, _, _, _, errParseName := parseName(testcase.element)
			_, _, errParseElement := parserConfiguration{"failing", expectedElements{}}.parseElement([]string{testcase.element})

			assertions := require.New(t)
			assertions.EqualError(errParseName, testcase.error)
			assertions.EqualError(errParseElement, errParseName.Error())
		})
	}
}

func TestParseNameNotRecursiveList(t *testing.T) {
	assertions := require.New(t)
	name, selector, isScalar, isRecursive, err := parseName("title")
	assertions.NoError(err)
	assertions.Equal("title", name)
	assertions.Equal("", selector)
	assertions.True(isScalar)
	assertions.False(isRecursive)
}

func TestParseNameListNoSelector(t *testing.T) {
	assertions := require.New(t)
	name, selector, isScalar, isRecursive, err := parseName("changes[]")
	assertions.NoError(err)
	assertions.Equal("changes", name)
	assertions.Equal("", selector)
	assertions.False(isScalar)
	assertions.False(isRecursive)
}

func TestParseNameListWithSelector(t *testing.T) {
	assertions := require.New(t)
	name, selector, isScalar, isRecursive, err := parseName("changes[2]")
	assertions.NoError(err)
	assertions.Equal("changes", name)
	assertions.Equal("2", selector)
	assertions.False(isScalar)
	assertions.False(isRecursive)
}

func TestParseNameRecursiveListNoSelector(t *testing.T) {
	assertions := require.New(t)
	name, selector, isScalar, isRecursive, err := parseName("changes[]/")
	assertions.NoError(err)
	assertions.Equal("changes", name)
	assertions.Equal("", selector)
	assertions.False(isScalar)
	assertions.True(isRecursive)
}
