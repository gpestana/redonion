package outputs

type Unit struct {
	Url     string
	Outputs []interface{}
}

type Output interface {
	Handle(Unit) error
}
