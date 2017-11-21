package processor

type Processor interface {
	InChannel() chan DataUnit
	Name() string
	Process()
	Stop() error
}

type DataUnit struct {
	Processor *Processor
	Output    string
}
