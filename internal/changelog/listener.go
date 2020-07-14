package changelog

type Listener interface {
	Enter(h Heading)
	Exit(h Heading)
}
