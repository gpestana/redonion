package processor

import (
	"log"
)

type TextProcessor struct {
	name        string
	next        *Processor
	inChannel   chan string
	outChannel  chan string
	inputLenght int
}

func NewTextProcessor(in chan string, out chan string, len int) TextProcessor {
	return TextProcessor{
		name:        "text",
		inChannel:   in,
		outChannel:  out,
		inputLenght: len,
	}
}

func (p TextProcessor) Process() {
	log.Println("TextProcessor.Process")
	for i := 0; i < p.inputLenght; i++ {
		msg := <-p.inChannel
		log.Println(len(msg))

		// pass to next processing unit
		p.outChannel <- msg
	}
}

func (p TextProcessor) Stop() error {
	log.Println("TextProcessor.Stop")
	return nil
}
