package outputs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	docType = "unit"
)

type EsOutput struct {
	user     string
	password string
	host     string
	index    string
}

type PostResponse struct {
	Index  string `json:"_index"`
	Type   string `json:"_type"`
	Id     string `json:"_id"`
	Result string `json:"result"`
}

func NewEs(u string, p string, h string, i string) (EsOutput, error) {
	return EsOutput{
		user:     u,
		password: p,
		host:     h,
		index:    i,
	}, nil
}

func (out *EsOutput) Handle(un Unit) error {
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
