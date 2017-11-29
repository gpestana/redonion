package processor

import (
	"github.com/google/uuid"
	"io"
	"log"
)

type Processor interface {
	InChannel() chan DataUnit
	Name() string
	Process()
}

type DataUnit struct {
	Processor *Processor
	Url       string
	Output    []byte
	Reader    io.Reader
}

func Name(n string) string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return n + "_" + uuid.String()
}
