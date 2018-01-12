package errors

import (
	"log"
)

const (
	FetcherError = "fetcherError"
	ProcessError = "processError"
	OutputError  = "outputError"
)

type Error struct {
	Type    string
	Message string
}

func (e Error) Print() {
	log.Println(e)
}
