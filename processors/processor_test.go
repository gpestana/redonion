package processor

import (
	"log"
	"os"
	"testing"
)

func TestProcessor_conf(t *testing.T) {
	c := `{
	"processors": [{
		"type": "image",
		"tensorflow_url": "http://localhost:8080/recognize"
	}]
}
`
	p := "./test_conf.json"
	defer os.Remove(p)
	f, err := os.Create(p)
	if err != nil {
		log.Fatal("config loading error")
	}

	_, err = f.Write([]byte(c))
	if err != nil {
		log.Fatal("config loading error")
	}

	cnf, err := ParseConfig(&p)
	if err != nil {
		t.Error(err)
	}

	if tp := cnf.Processors[0].Type; tp != "image" {
		t.Error("configuration parsing error: image !=" + tp)
	}
	if tf := cnf.Processors[0].TFUrl; tf != "http://localhost:8080/recognize" {
		t.Error("configuration parsing error: http://localhost:8080/recognize != " + tf)
	}
}
