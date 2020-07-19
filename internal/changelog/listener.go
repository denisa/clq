package changelog

// a Listener is notified as a heading is visited
type Listener interface {
	// Enter is called when a heading is first met.
	Enter(h Heading)
	// Exit is called when all of a heading’s children have been visited.
	Exit(h Heading)
}
