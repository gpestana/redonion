package output

type Conf struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Index    string `json:"index"`
}

type Unit struct {
	Url     string
	Outputs []interface{}
}

type Output interface {
	Start()
	Handle(Unit) error
	Results() ([]byte, error)
}
