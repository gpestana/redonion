package processor

import (
	"bytes"
	"encoding/json"
	"github.com/gpestana/redonion/tor"
	"github.com/xiam/exif"
	"golang.org/x/net/html"
	"io"
	"log"
	"strings"
)

type ImageProcessor struct {
	name        string
	inChannel   chan DataUnit
	outChannel  chan DataUnit
	inputLength int
	images      []Image
}

type Image struct {
	RootUrl       string
	Url           string
	ProcessorName string
	Exif          map[string]string
	Recon         []string
	Errors        []string
}

func NewImageProcessor(in chan DataUnit, out chan DataUnit, len int) ImageProcessor {
	return ImageProcessor{
		name:        Name("image"),
		inChannel:   in,
		outChannel:  out,
		inputLength: len,
	}
}

func (img Image) Json() ([]byte, error) {
	json, err := json.Marshal(img)
	return json, err
}

func (p ImageProcessor) Process() {
	for j := 0; j < p.inputLength; j++ {
		du := DataUnit{}
		du = <-p.inChannel

		imgUrls := images(du.Html)
		for _, url := range imgUrls {
			errs := []string{}
			curl := canonicalUrl(du.Url, url)
			imgData, err := tor.Get(curl)
			if err != nil {
				errs = append(errs, err.Error())
			}
			meta, err := metadata(imgData)
			if err != nil {
				errs = append(errs, err.Error())
			}

			i := Image{
				RootUrl:       du.Url,
				Url:           curl,
				ProcessorName: p.name,
				Exif:          meta,
				Recon:         []string{},
				Errors:        errs,
			}
			du.Outputs = append(du.Outputs, i)
		}
		p.outChannel <- du
	}
}

func (p ImageProcessor) Name() string {
	return p.name
}

func (p ImageProcessor) InChannel() chan DataUnit {
	return p.inChannel
}

//gets all image urls from HTML
func images(b []byte) []string {
	urls := []string{}
	r := bytes.NewReader(b)
	tz := html.NewTokenizer(r)
	for {
		tok := tz.Next()
		switch {
		case tok == html.ErrorToken:
			return urls
		case tok == html.StartTagToken:
			t := tz.Token()
			if t.Data == "img" {
				for _, a := range t.Attr {
					if a.Key == "src" {
						urls = append(urls, a.Val)
					}
				}
			}
		case tok == html.SelfClosingTagToken:
			t := tz.Token()
			if t.Data == "img" {
				for _, a := range t.Attr {
					if a.Key == "src" {
						urls = append(urls, a.Val)
					}
				}
			}
		}
	}
}

//gets image metadata if possible
func metadata(data []byte) (map[string]string, error) {
	r := exif.New()
	buf := bytes.NewBuffer(data)
	_, err := io.Copy(r, buf)
	if err.Error() != "Found EXIF header. OK to call Parse." {
		return nil, err
	}
	err = r.Parse()
	if err != nil {
		return nil, err
	}
	return r.Tags, nil
}

//gets recognition info about image
func (img *Image) recon() {
	log.Println("Image.Recon")
}

func canonicalUrl(b string, u string) string {
	if strings.HasPrefix(u, "http") || strings.HasPrefix(u, "www") {
		return u
	}
	b = strings.TrimSuffix(b, "/")
	u = strings.TrimPrefix(u, ".")
	u = strings.TrimPrefix(u, "/")
	return b + "/" + u
}
