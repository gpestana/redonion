package tor

import (
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	torProxy = "127.0.0.1:9150"
	timeout  = 15
)

func Get(u string) ([]byte, error) {
	proxyURL, err := url.Parse("socks5://" + torProxy)
	if err != nil {
		log.Println("Failed to parse proxy URL: " + u)
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Println("Failed to obtain proxy dialer ", err)
		return nil, err
	}

	t := &http.Transport{Dial: dialer.Dial}
	c := &http.Client{
		Transport: t,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	r, err := c.Get(u)
	if err != nil {
		log.Println("Failed to issue GET request: ", err)
		return nil, nil
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read the body ", err)
	}
	return b, nil
}
