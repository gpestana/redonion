package main

import (
	"log"
	"os"
	"testing"
)

func TestConfigOutputES(t *testing.T) {
	p := "./test_c.json"
	c := `{
  "outputs": [{
      "type": "elasticsearch",
      "user": "user",
      "password": "pass",
      "host": "localhost:9200",
      "index": "index-01"
    }
  ]
}	
`
	defer os.Remove(p)
	createConfigFile(p, c)
	conf, err := config(&p)
	if err != nil {
		t.Error(err.Error())
	}

	if conf.Outputs[0].Type != "elasticsearch" {
		t.Error("elasticsearch type output not parsed correctly")
	}

	if conf.Outputs[0].User != "user" {
		t.Error("elasticsearch user output not parsed correctly")
	}

	if conf.Outputs[0].Password != "pass" {
		t.Error("elasticsearch pass output not parsed correctly")
	}

	if conf.Outputs[0].Host != "localhost:9200" {
		t.Error("elasticsearch type host not parsed correctly")
	}

	if conf.Outputs[0].Index != "index-01" {
		t.Error("elasticsearch type index not parsed correctly")
	}

}

func createConfigFile(p string, conf string) {
	f, err := os.Create(p)
	if err != nil {
		log.Fatal("test setup failed" + err.Error())
	}
	_, err = f.Write([]byte(conf))

}
