package processor

import (
	"bytes"
)

type Processor struct {
	name      string
	next      *Processor
	InChannel chan bytes.Buffer
}

func New(n string) (Processor, error) {
	return Processor{
		name: n,
	}, nil
}

func (p Processor) Process() {}
func (p Processor) stop()    {}
func start(p []string)       {}
