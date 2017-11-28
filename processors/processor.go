package processor

import (
	"github.com/google/uuid"
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
}

func Name(n string) string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return n + "_" + uuid.String()
}
