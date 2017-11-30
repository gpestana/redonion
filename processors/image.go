package processor

import (
	"golang.org/x/net/html"
	//"html"
	"io"
	"log"
)

type ImageProcessor struct {
	name        string
	inChannel   chan DataUnit
	outChannel  chan DataUnit
	inputLength int
	images      []Image
}

type Image struct {
	url      string
	metadata []string
	recon    []string
}

func NewImageProcessor(in chan DataUnit, out chan DataUnit, len int) ImageProcessor {
	return ImageProcessor{
		name:        Name("image"),
		inChannel:   in,
		outChannel:  out,
		inputLength: len,
		images:      []Image{},
	}
}

func (p *ImageProcessor) Process() {
	log.Println("ImageProcessor.Process")
	for j := 0; j < p.inputLength; j++ {
		du := DataUnit{}
		du = <-p.inChannel

		imgs := images(du.Reader)
		for _, url := range imgs {
			i := Image{
				url:      url,
				metadata: []string{},
				recon:    []string{},
			}
			i.Metadata()
			i.Recon()
			p.images = append(p.images, i)
		}

		p.outChannel <- du
	}
}

func (p *ImageProcessor) Name() string {
	return p.name
}

func (p *ImageProcessor) InChannel() chan DataUnit {
	return p.inChannel
}

//gets all image urls from HTML
func images(r io.Reader) []string {
	urls := []string{}

	tz := html.NewTokenizer(r)
	for {
		tok := tz.Next()
		switch {
		case tok == html.ErrorToken:
			return urls
		case tok == html.StartTagToken:
			t := tz.Token()
			if t.Data == "img" {
				urls = append(urls, "img")
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
func (img *Image) Metadata() {
	log.Println("Image.Metadata")
}

//gets recognition info about image
func (img *Image) Recon() {
	log.Println("Image.Recon")
}
