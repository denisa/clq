package changelog

import (
	"log"
)

type Logger struct{}

func (l Logger) Enter(h Heading) {
	log.Printf("Entering %v", h)
}

func (l Logger) Exit(h Heading) {
	log.Printf("Exiting %v", h)
}
