package fetcher

import (
	"bufio"
	"fmt"
	"github.com/gpestana/redonion/processors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Fetcher struct {
	urls       []string
	proxy      string
	listPath   string
	timeout    int
	processors []processor.Processor
}

func New(urlsIn *string, proxy *string, listPath *string, timeout *int, proc []string) (Fetcher, error) {
	//do input verification
	var urls []string

	if *listPath != "" {
		urlsPath, err := parseListURL(*listPath)
		if err != nil {
			return Fetcher{}, err
		}
		urls = urlsPath
	} else {
		urls = strings.Split(*urlsIn, ",")
	}

	processors := []processor.Processor{}
	for _, n := range proc {
		p, err := processor.New(n)
		if err != nil {
			log.Fatal(err)
		}
		processors = append(processors, p)
	}

	return Fetcher{
		urls:       urls,
		proxy:      *proxy,
		listPath:   *listPath,
		timeout:    *timeout,
		processors: processors,
	}, nil
}

func (f *Fetcher) Start() []string {
	channel := make(chan string)
	var results []string

	for _, u := range f.urls {
		go func(u string) {
			r, _ := f.request(u)
			channel <- r
		}(u)
	}

	for range f.urls {
		r := <-channel
		results = append(results, r)
	}
	close(channel)
	return results
}

func parseListURL(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var urls []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		urls = append(urls, s.Text())
	}
	return urls, s.Err()
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
