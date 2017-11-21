package main

import (
	"bufio"
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"github.com/gpestana/redonion/outputs"
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

	text1Chn := make(chan processor.DataUnit, len(ulist))
	text2Chn := make(chan processor.DataUnit, len(ulist))
	chs := []chan processor.DataUnit{text1Chn, text2Chn}

	outputChn := make(chan processor.DataUnit, len(ulist)*len(chs))
	output := output.NewStdout(outputChn, len(ulist))

	processors := []processor.Processor{
		processor.NewTextProcessor(text1Chn, outputChn, len(ulist)),
		processor.NewTextProcessor(text2Chn, outputChn, len(ulist)),
	}

	fetcher, err := fetcher.New(ulist, proxy, timeout, processors)
	if err != nil {
		log.Fatal(err)
	}

	fetcher.Start()

	// start processors
	for _, p := range processors {
		p.Process()
	}

	// run output
	for i := 0; i < len(ulist); i++ {
		output.Now()
	}

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
