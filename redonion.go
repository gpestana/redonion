package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"github.com/gpestana/redonion/processors"
	"log"
	"os"
	"strings"
)

// data structure in which outputs will be rendered
type Output struct {
	Url     string
	Outputs []interface{}
}

func main() {

	urls := flag.String("urls", "http://127.0.0.1", "list of addresses to scan (separated by comma)")
	list := flag.String("list", "", "path for list of addresses to scan")
	flag.Parse()

	ulist, err := parseUrls(urls, list)
	if err != nil {
		log.Fatal(err)
	}

	imgChn := make(chan processor.DataUnit, len(ulist))
	hiddenChn := make(chan processor.DataUnit, len(ulist))
	chs := []chan processor.DataUnit{hiddenChn, imgChn}

	outputChn := make(chan processor.DataUnit, len(ulist)*len(chs))

	processors := []processor.Processor{
		processor.NewHiddenProcessor(hiddenChn, outputChn, len(ulist)),
		processor.NewImageProcessor(imgChn, outputChn, len(ulist)),
	}

	fetcher, err := fetcher.New(ulist, processors)
	if err != nil {
		log.Fatal(err)
	}

	// start fetcher
	fetcher.Start()

	// spin up all processors
	for _, p := range processors {
		p.Process()
	}

	// get all results from output channel
	results := []Output{}
	sizeRes := len(ulist) * len(chs)
	for i := 0; i < sizeRes; i++ {
		du := <-outputChn
		r := Output{
			Url:     du.Url,
			Outputs: du.Outputs,
		}
		results = append(results, r)
	}

	jres, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(jres)

	closeChannels(chs, outputChn)
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

func closeChannels(chs []chan processor.DataUnit, outputCh chan processor.DataUnit) {
	for _, ch := range chs {
		close(ch)
	}
	close(outputCh)
}
