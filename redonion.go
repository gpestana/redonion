package main

import (
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"github.com/gpestana/redonion/processors"
	"log"
)

func main() {

	urls := flag.String("urls", "http://127.0.0.1", "list of addresses to scan (separated by comma)")
	list := flag.String("list", "", "path for list of addresses to scan")
	timeout := flag.Int("timeout", 15, "requests timeout (seconds)")
	proxy := flag.String("proxy", "127.0.0.1:9150", "url of tor proxy")
	flag.Parse()

	// initializes the pipeline
	inputChn := make(chan string)
	textChn := make(chan string)

	processors := []processor.Processor{
		processor.NewTextProcessor(inputChn, textChn, 2),
	}

	fetcher, err := fetcher.New(urls, proxy, list, timeout, inputChn)
	if err != nil {
		log.Fatal(err)
	}

	fetcher.Start()

	for _, p := range processors {
		p.Process()
	}
}
