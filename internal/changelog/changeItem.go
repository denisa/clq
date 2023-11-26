package changelog

import "fmt"

// A single change item description
type ChangeItem struct {
	heading
}

func (h HeadingsFactory) newChangeItem(title string) (Heading, error) {
	if title == "" {
		return nil, fmt.Errorf("Validation error: change description cannot stay empty")
	}
	return ChangeItem{heading{title: title, kind: ChangeDescription}}, nil
}

func (h ChangeItem) DisplayTitle() string {
	return h.Title()
}

func (h ChangeItem) Title() string {
	return h.title
}

func (h ChangeItem) Kind() HeadingKind {
	return h.kind
}

func (h ChangeItem) String() string {
	return asPath(h.title)
}
