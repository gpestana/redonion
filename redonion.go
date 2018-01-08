package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"github.com/gpestana/redonion/output"
	"github.com/gpestana/redonion/processors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	urls := flag.String("urls", "http://127.0.0.1", "list of addresses to scan (separated by comma)")
	list := flag.String("list", "", "path for list of addresses to scan")
	c := flag.String("config", "", "path for configuration file")
	flag.Parse()

	cnf, err := config(c)
	procCnf, err := processor.ParseConfig(c)
	if err != nil {
		log.Fatal("Could not parse configuration file " + err.Error())
	}

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
		processor.NewImageProcessor(imgChn, outputChn, len(ulist), procCnf),
	}

	sizeResults := len(ulist) * len(chs)
	outputs := []output.Output{
		output.NewEs(cnf.Outputs, outputChn, sizeResults),
	}

	fetcher, err := fetcher.New(ulist, processors)
	if err != nil {
		log.Fatal(err)
	}
	fetcher.Start()

	for _, p := range processors {
		p.Process()
	}

	for _, o := range outputs {
		o.Start()
	}

	// write outputs to stdout
	for _, o := range outputs {
		jres, err := o.Results()
		if err != nil {
			log.Println(err)
		} else {
			os.Stdout.Write(jres)
		}
	}

	closeChannels(chs, outputChn)
}

type processorsC struct {
	Type  string `json:"type"`
	TFUrl string `json:"tensorflow_url"`
}

type Config struct {
	Outputs    []output.Conf `json:"outputs"`
	Processors []processorsC `json:"processors"`
}

func config(p *string) (Config, error) {
	cf, err := ioutil.ReadFile(*p)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = json.Unmarshal([]byte(cf), &config)
	if err != nil {
		return Config{}, nil
	}
	return config, nil
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
