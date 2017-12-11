//processor that gathers all hidden onion URLs from the fetched content
package processor

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
)

type HiddenProcessor struct {
	name           string
	inChannel      chan DataUnit
	outChannel     chan DataUnit
	inputLength    int
	HiddenServices []HiddenService
}

type HiddenService struct {
	RootUrl       string
	ProcessorName string
	Url           string
}

func NewHiddenProcessor(in chan DataUnit, out chan DataUnit, len int) HiddenProcessor {
	return HiddenProcessor{
		name:        Name("image"),
		inChannel:   in,
		outChannel:  out,
		inputLength: len,
	}
}

func (p HiddenProcessor) InChannel() chan DataUnit {
	return p.inChannel
}

func (p HiddenProcessor) Process() {
	for j := 0; j < p.inputLength; j++ {
		du := DataUnit{}
		du = <-p.inChannel
		hs := hiddenUrls(du.Html)

		for _, u := range hs {
			h := HiddenService{
				RootUrl:       du.Url,
				ProcessorName: p.name,
				Url:           u,
			}
			du.Outputs = append(du.Outputs, h)
		}
		p.outChannel <- du
	}
}

//get all hidden services from HTML
func hiddenUrls(b []byte) []string {
	urls := []string{}
	r := bytes.NewReader(b)
	tz := html.NewTokenizer(r)
	for {
		tok := tz.Next()
		switch {
		case tok == html.ErrorToken:
			return urls
		case tok == html.StartTagToken:
			t := tz.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if strings.Contains(a.Val, ".onion") {
							urls = append(urls, a.Val)
						}
					}
				}
			}
		}
	}
}
