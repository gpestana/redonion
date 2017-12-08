package fetcher

import (
	"github.com/gpestana/redonion/processors"
	"github.com/gpestana/redonion/tor"
	"log"
)

type Fetcher struct {
	urls       []string
	processors []processor.Processor
}

func New(urls []string, pr []processor.Processor) (Fetcher, error) {
	//do input verification
	return Fetcher{
		urls:       urls,
		processors: pr,
	}, nil
}

func (f *Fetcher) Start() {
	for _, u := range f.urls {
		log.Println("Fetcher: spinning goroutine " + u)
		go func(u string) {
			//out := []processor.Output{}
			var out []interface{}
			b, _ := tor.Get(u)
			// fan-out result from fetcher to all registerd processors
			for _, p := range f.processors {
				du := processor.DataUnit{&p, u, b, out}
				p.InChannel() <- du
			}
		}(u)
	}
}
