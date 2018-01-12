package fetcher

import (
	"github.com/gpestana/redonion/errors"
	"github.com/gpestana/redonion/processors"
	"github.com/gpestana/redonion/tor"
)

type Fetcher struct {
	urls       []string
	processors []processor.Processor
	errors     *[]errors.Error
}

func New(urls []string, pr []processor.Processor, err *[]errors.Error) (Fetcher, error) {
	return Fetcher{
		urls:       urls,
		processors: pr,
		errors:     err,
	}, nil
}

func (f *Fetcher) Start() {
	for _, u := range f.urls {
		go func(u string) {
			var out []interface{}
			b, err := tor.Get(u)
			if err != nil {
				*f.errors = append(*f.errors, errors.Error{errors.FetcherError, err.Error()})
			}

			for _, p := range f.processors {
				du := processor.DataUnit{&p, u, b, out}
				p.InChannel() <- du
			}
		}(u)
	}
}
