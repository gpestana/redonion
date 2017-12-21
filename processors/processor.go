package processor

import (
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
)

type processorsC struct {
	Type  string `json:"type"`
	TFUrl string `json:'tensorflow_url'`
}

type Config struct {
	Processors []processorsC `json:"processors"`
}

type Processor interface {
	InChannel() chan DataUnit
	Process()
}

func Name(n string) string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return n + "_" + uuid.String()
}

type DataUnit struct {
	Processor *Processor
	Url       string
	Html      []byte
	Outputs   []interface{}
}

type Output interface {
	Json() ([]byte, error)
}

func ParseConfig(p *string) (Config, error) {
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
