package processor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gpestana/redonion/tor"
	"github.com/xiam/exif"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
)

type ImageProcessor struct {
	name        string
	inChannel   chan DataUnit
	outChannel  chan DataUnit
	inputLength int
	images      []Image
	tfUrl       string
}

type Image struct {
	RootUrl       string
	Url           string
	ProcessorName string
	Exif          map[string]string
	Recon         ReconResults
	Errors        []string
	Timestamp     time.Time
}

func NewImageProcessor(in chan DataUnit, out chan DataUnit, len int, cnf Config) ImageProcessor {
	var tfUrl string
	for _, c := range cnf.Processors {
		if c.Type == "image" {
			tfUrl = c.TFUrl
			break
		}
	}

	return ImageProcessor{
		name:        Name("image"),
		inChannel:   in,
		outChannel:  out,
		inputLength: len,
		tfUrl:       tfUrl,
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

			// TODO: refactor to struct method?
			meta, err := metadata(imgData)
			if err != nil {
				errs = append(errs, err.Error())
			}

			recons, err := recognition(imgData, p.tfUrl)
			if err != nil {
				errs = append(errs, err.Error())
			}

			i := Image{
				RootUrl:       du.Url,
				Url:           curl,
				ProcessorName: p.name,
				Exif:          meta,
				Recon:         recons,
				Errors:        errs,
				Timestamp:     time.Now(),
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

type ReconResults struct {
	Results []Recon
}

type Recon struct {
	Label       string
	Probability float32
}

func recognition(d []byte, tfUrl string) (ReconResults, error) {
	if tfUrl == "" {
		return ReconResults{}, errors.New("Tensorflow URL is not defined")
	}
	data := bytes.NewBuffer(d)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name=image; filename="from-buffer"`))
	fw, err := w.CreatePart(h)
	if err != nil {
		return ReconResults{}, err
	}
	if _, err = io.Copy(fw, data); err != nil {
		return ReconResults{}, err
	}

	w.Close()
	req, err := http.NewRequest("POST", tfUrl, &b)
	if err != nil {
		return ReconResults{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	cli := &http.Client{}
	r, err := cli.Do(req)
	if err != nil {
		return ReconResults{}, err
	}

	if r == nil {
		return ReconResults{}, errors.New("Recon: HTTP response is empty")
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	type tfResults struct {
		Filename string
		Labels   []Recon
		Err      string `json:"error"`
	}

	res := tfResults{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return ReconResults{}, err
	}
	if recErr := res.Err; recErr != "" {
		return ReconResults{}, errors.New("Recon: " + recErr)
	}

	recons := ReconResults{res.Labels}
	return recons, nil
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
