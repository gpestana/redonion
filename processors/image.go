package processor

import (
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

		imgs := images(du.Url)
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

func images(url string) []string {
	urls := []string{}
	log.Println("Image.images")
	return urls
}

func (img *Image) Metadata() {
	log.Println("Image.Metadata")
}

func (img *Image) Recon() {
	log.Println("Image.Recon")
}
