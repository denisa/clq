package changelog

import (
	"fmt"
)

// recorder keeps a logs of events (entering ot exiting a header) for use in unit testings.
type recorder struct {
	events []string
}

func (r *recorder) Enter(h Heading) {
	r.events = append(r.events, fmt.Sprintf("Enter %v", h))
}

func (r *recorder) Exit(h Heading) {
	r.events = append(r.events, fmt.Sprintf("Exit %v", h))
}
