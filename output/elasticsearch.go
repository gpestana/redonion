package output

import (
	"bytes"
	"encoding/json"
	"github.com/gpestana/redonion/processors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	docType = "unit"
)

type EsOutput struct {
	user        string
	password    string
	host        string
	index       string
	inputChan   chan processor.DataUnit
	sizeResults int
	results     []Unit
}

type PostResponse struct {
	Index  string `json:"_index"`
	Type   string `json:"_type"`
	Id     string `json:"_id"`
	Result string `json:"result"`
}

func NewEs(confs []Conf, in chan processor.DataUnit, sizeRes int) EsOutput {
	for _, c := range confs {
		if c.Type == "elasticsearch" {
			return EsOutput{
				user:        c.User,
				password:    c.Password,
				host:        c.Host,
				index:       c.Index,
				inputChan:   in,
				sizeResults: sizeRes,
			}
		}
	}
	return EsOutput{}
}

func (out EsOutput) Start() {
	results := []Unit{}
	for i := 0; i < out.sizeResults; i++ {
		du := <-out.inputChan
		r := Unit{
			Url:     du.Url,
			Outputs: du.Outputs,
		}
		err := out.Handle(r)
		if err != nil {
			// error handling
			log.Println(err)
		}
		results = append(results, r)
	}
	out.results = results
}

func (out EsOutput) Handle(un Unit) error {
	cli := http.Client{}
	u := out.host + "/" + out.index + "/" + docType

	for _, out := range un.Outputs {
		b, err := json.Marshal(out)
		if err != nil {
			return err
		}
		br := bytes.NewReader(b)
		req, err := http.NewRequest("POST", u, br)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		r, err := cli.Do(req)
		if err != nil {
			return err
		}
		defer r.Body.Close()
		rb, _ := ioutil.ReadAll(r.Body)
		resp := PostResponse{}
		_ = json.Unmarshal(rb, &resp)
		if resp.Result != "created" {
			log.Println("Outputs.Elasticsearch.Handle:: " + resp.Result)
		}
	}

	return nil
}

func (out EsOutput) Results() ([]byte, error) {
	jres, err := json.Marshal(out.results)
	if err != nil {
		return []byte{}, err
	}
	return jres, nil
}
