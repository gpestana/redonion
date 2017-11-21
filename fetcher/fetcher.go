package fetcher

import (
	"fmt"
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
	outChannel chan string
}

func New(urls []string, proxy *string, timeout *int, out chan string) (Fetcher, error) {
	//do input verification
	return Fetcher{
		urls:       urls,
		proxy:      *proxy,
		timeout:    *timeout,
		outChannel: out,
	}, nil
}

func (f *Fetcher) Start() {
	log.Println("Fetcher.Start")
	for _, u := range f.urls {
		log.Println("Fetcher.Start: spinning new goroutine " + u)
		go func(u string) {
			r, _ := f.request(u)
			f.outChannel <- r
		}(u)
	}
}

func (f *Fetcher) request(u string) (string, error) {
	proxyURL, err := url.Parse("socks5://" + f.proxy)
	if err != nil {
		fmt.Println("Failed to parse proxy URL: " + u)
		return "", err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		fmt.Println("Failed to obtain proxy dialer ", err)
		return "", err
	}

	t := &http.Transport{Dial: dialer.Dial}
	c := &http.Client{
		Transport: t,
		Timeout:   time.Duration(f.timeout) * time.Second,
	}
	r, err := c.Get(u)
	if err != nil {
		fmt.Println("Failed to issue GET request: ", err)
		return err.Error(), nil
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failed to read the body ", err)
	}
	return string(b[:]), nil
}
