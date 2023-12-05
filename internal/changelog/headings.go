package changelog

// HeadingKind is the type for the multiple sections.
type HeadingKind int

const (
	IntroductionHeading HeadingKind = iota
	ReleaseHeading
	ChangeHeading
	ChangeDescription
)

// A Heading is the interface common to every sections.
type Heading interface {
	// The Title of the section
	Title() string
	// DisplayTitle is the title to be displayed
	DisplayTitle() string
	// Kind is the HeadingKind for the section
	Kind() HeadingKind
	String() string
}

type heading struct {
	title string
	kind  HeadingKind
}

func asPath(name string) string {
	return "{" + name + "}"
}
