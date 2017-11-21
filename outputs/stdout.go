package output

import (
	"github.com/gpestana/redonion/processors"
	"log"
)

type Stdout struct {
	inChannel    chan processor.DataUnit
	outputLength int
}

func NewStdout(chn chan processor.DataUnit, len int) Stdout {
	return Stdout{
		inChannel:    chn,
		outputLength: len,
	}
}

func (o Stdout) Now() {
	for i := 0; i < o.outputLength; i++ {
		du := <-o.inChannel
		log.Println(len(du.Output))
	}
}
