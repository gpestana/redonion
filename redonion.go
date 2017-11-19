package main

import (
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"log"
)

func main() {

	urls := flag.String("urls", "http://127.0.0.1", "list of addresses to scan (separated by comma)")
	list := flag.String("list", "", "path for list of addresses to scan")
	timeout := flag.Int("timeout", 15, "requests timeout (seconds)")
	proxy := flag.String("proxy", "127.0.0.1:9150", "url of tor proxy")
	flag.Parse()

	processors := []string{"image", "text"}

	fetcher, err := fetcher.New(urls, proxy, list, timeout, processors)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fetcher)
	//results := fetcher.Start()
	//log.Println(results)
}
