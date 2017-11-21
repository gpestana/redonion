package main

import (
	"bufio"
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"github.com/gpestana/redonion/processors"
	"log"
	"os"
	"strings"
)

func main() {

	urls := flag.String("urls", "http://127.0.0.1", "list of addresses to scan (separated by comma)")
	list := flag.String("list", "", "path for list of addresses to scan")
	timeout := flag.Int("timeout", 15, "requests timeout (seconds)")
	proxy := flag.String("proxy", "127.0.0.1:9150", "url of tor proxy")
	flag.Parse()

	ulist, err := parseUrls(urls, list)
	if err != nil {
		log.Fatal(err)
	}

	// init pipeline
	inputChn := make(chan string, len(ulist))
	textChn := make(chan string, len(ulist))
	text2Chn := make(chan string, len(ulist))
	chs := []chan string{inputChn, textChn, text2Chn}

	processors := []processor.Processor{
		processor.NewTextProcessor(inputChn, textChn, len(ulist)),
		processor.NewTextProcessor(textChn, text2Chn, len(ulist)),
	}

	fetcher, err := fetcher.New(ulist, proxy, timeout, inputChn)
	if err != nil {
		log.Fatal(err)
	}

	// start pipeline
	fetcher.Start()
	for _, p := range processors {
		p.Process()
	}

	// cleanup
	closeChannels(chs)
}

func parseUrls(urlsIn *string, list *string) ([]string, error) {
	var urls []string
	if *list != "" {
		urlsPath, err := parseListURL(*list)
		if err != nil {
			return nil, err
		}
		urls = urlsPath
	} else {
		urls = strings.Split(*urlsIn, ",")
	}
	return urls, nil
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

func closeChannels(chs []chan string) {
	for _, ch := range chs {
		close(ch)
	}
}
