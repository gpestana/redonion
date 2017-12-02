package output

import (
	"encoding/json"
	"github.com/gpestana/redonion/processors"
)

type Stdout struct {
	inChannel    chan processor.DataUnit
	outputLength int
	Results      []StdoutResult
}

type StdoutResult struct {
	Url           string
	ProcessorName string
	Output        []byte
	OutputString  string
}

func NewStdout(chn chan processor.DataUnit, len int) Stdout {
	return Stdout{
		inChannel:    chn,
		outputLength: len,
		Results:      []StdoutResult{},
	}
}

func (o *Stdout) Run() {
	for i := 0; i < o.outputLength; i++ {
		du := <-o.inChannel
		pr := *du.Processor
		r := StdoutResult{
			Url:           du.Url,
			ProcessorName: pr.Name(),
			Output:        du.Output,
			OutputString:  string(du.Output),
		}
		o.Results = append(o.Results, r)
	}
}

func (o Stdout) Result() ([]byte, error) {
	json, err := json.Marshal(o.Results)
	return json, err
}
