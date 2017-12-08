package processor

type TextProcessor struct {
	name        string
	inChannel   chan DataUnit
	outChannel  chan DataUnit
	inputLenght int
}

type Text struct {
	Url           string
	ProcessorName string
	Text          string
	Error         error
}

func NewTextProcessor(in chan DataUnit, out chan DataUnit, len int) TextProcessor {
	return TextProcessor{
		name:        Name("text"),
		inChannel:   in,
		outChannel:  out,
		inputLenght: len,
	}
}

func (p TextProcessor) Process() {
	for i := 0; i < p.inputLenght; i++ {
		du := DataUnit{}
		du = <-p.inChannel
		t := Text{
			Url:           du.Url,
			ProcessorName: p.name,
			Text:          string(du.Html),
		}
		du.Outputs = append(du.Outputs, t)
		p.outChannel <- du
	}
}

func (p TextProcessor) Name() string {
	return p.name
}

func (p TextProcessor) InChannel() chan DataUnit {
	return p.inChannel
}
