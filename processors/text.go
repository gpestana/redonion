package processor

import (
	"github.com/google/uuid"
	"log"
)

type TextProcessor struct {
	name        string
	inChannel   chan DataUnit
	outChannel  chan DataUnit
	inputLenght int
}

func NewTextProcessor(in chan DataUnit, out chan DataUnit, len int) TextProcessor {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	name := "text" + uuid.String()
	return TextProcessor{
		name:        name,
		inChannel:   in,
		outChannel:  out,
		inputLenght: len,
	}
}

func (p TextProcessor) Process() {
	log.Println("TextProcessor.Process")
	for i := 0; i < p.inputLenght; i++ {
		du := DataUnit{}
		du = <-p.inChannel

		// pass to next processing unit
		p.outChannel <- du
	}
}

func (p TextProcessor) Stop() error {
	log.Println("TextProcessor.Stop")
	return nil
}

func (p TextProcessor) Name() string {
	return p.name
}
func (p TextProcessor) InChannel() chan DataUnit {
	return p.inChannel
}
