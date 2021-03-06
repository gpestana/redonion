package tor

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	torProxy = "127.0.0.1:9150"
	timeout  = 30
)

func Get(u string) ([]byte, error) {
	proxyURL, err := url.Parse("socks5://" + torProxy)
	if err != nil {
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return nil, err
	}

	t := &http.Transport{Dial: dialer.Dial}
	c := &http.Client{
		Transport: t,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	r, err := c.Get(u)
	//log.Println(err)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
	}
	return b, nil
}
