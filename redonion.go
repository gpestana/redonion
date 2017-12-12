package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/gpestana/redonion/fetcher"
	"github.com/gpestana/redonion/outputs"
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

	// elasticsearch output
	es := outputs.EsOutput{}
	for _, c := range cnf.Outputs {
		if c.Type == "elasticsearch" {
			es, _ = outputs.NewEs(c.User, c.Password, c.Host, c.Index)
			break
		}
	}

	results := []outputs.Unit{}
	sizeRes := len(ulist) * len(chs)
	for i := 0; i < sizeRes; i++ {
		du := <-outputChn
		r := outputs.Unit{
			Url:     du.Url,
			Outputs: du.Outputs,
		}
		err := es.Handle(r)
		if err != nil {
			log.Println(err)
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

type outputConf struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Index    string `json:"index"`
}

type Config struct {
	Outputs []outputConf `json:"outputs"`
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
