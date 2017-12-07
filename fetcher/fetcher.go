package fetcher

import (
	"fmt"
	"github.com/gpestana/redonion/processors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Fetcher struct {
	urls       []string
	proxy      string
	timeout    int
	processors []processor.Processor
}

func New(urls []string, proxy *string, timeout *int, pr []processor.Processor) (Fetcher, error) {
	//do input verification
	return Fetcher{
		urls:       urls,
		proxy:      *proxy,
		timeout:    *timeout,
		processors: pr,
	}, nil
}

func (f *Fetcher) Start() {
	log.Println("Fetcher.Start")
	for _, u := range f.urls {
		log.Println("Fetcher.Start: spinning new goroutine " + u)
		go func(u string) {
			//out := []processor.Output{}
			var out []interface{}
			b, _ := f.request(u)
			// fan-out result from fetcher to all registerd processors
			for _, p := range f.processors {
				du := processor.DataUnit{&p, u, b, out}
				p.InChannel() <- du
			}
		}(u)
	}
}

func (f *Fetcher) request(u string) ([]byte, error) {
	proxyURL, err := url.Parse("socks5://" + f.proxy)
	if err != nil {
		fmt.Println("Failed to parse proxy URL: " + u)
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		fmt.Println("Failed to obtain proxy dialer ", err)
		return nil, err
	}

	t := &http.Transport{Dial: dialer.Dial}
	c := &http.Client{
		Transport: t,
		Timeout:   time.Duration(f.timeout) * time.Second,
	}
	r, err := c.Get(u)
	if err != nil {
		fmt.Println("Failed to issue GET request: ", err)
		return nil, nil
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failed to read the body ", err)
	}
	return b, nil
}
