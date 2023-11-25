package changelog

import (
	"fmt"
)

type HeadingsFactory struct {
	changeKind *ChangeKind
}

func NewHeadingFactory(changeKind *ChangeKind) HeadingsFactory {
	return HeadingsFactory{changeKind: changeKind}
}

// NewHeading is the factory method that, given a kind and a title, returns the appropriate Heading.
func (h HeadingsFactory) NewHeading(kind HeadingKind, title string) (Heading, error) {
	switch kind {
	case IntroductionHeading:
		return h.newIntroduction(title)
	case ReleaseHeading:
		return h.newRelease(title)
	case ChangeHeading:
		return h.newChange(title)
	case ChangeDescription:
		return h.newChangeItem(title)
	}
	return nil, fmt.Errorf("Unknown heading kind %v", kind)
}
