package processor

type Processor interface {
	Process()
	Stop() error
}
