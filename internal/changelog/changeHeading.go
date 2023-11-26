package changelog

import "fmt"

// Level 3, change groups
type Change struct {
	heading
	emoji string
}

func (h HeadingsFactory) newChange(title string) (Heading, error) {
	if title == "" {
		return nil, fmt.Errorf("Validation error: change cannot stay empty")
	}

	if emoji, err := h.changeKind.emojiFor(title); err != nil {
		return nil, err
	} else {
		return Change{heading{title: title, kind: ChangeHeading}, emoji}, nil
	}
}

func (h Change) DisplayTitle() string {
	if h.emoji == "" {
		return h.title
	}
	return h.emoji + " " + h.title
}

func (h Change) Title() string {
	return h.title
}

func (h Change) Kind() HeadingKind {
	return h.kind
}

func (h Change) String() string {
	return asPath(h.title)
}
