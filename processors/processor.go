package processor

import (
	"io"
)

type Processor interface {
	InChannel() chan DataUnit
	Name() string
	Process()
	Stop() error
}

type DataUnit struct {
	Processor *Processor
	Url       string
	Output    []byte
	Reader    io.Reader
}
