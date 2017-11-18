package main

import (
  "fmt"
	"log"
  "os/exec"
)

func main() {
	fmt.Println("starting")
	out, err := exec.Command(
		"./bin/onionscan",
		"--torProxyAddress=127.0.0.1:9150",
		"http://dmzwvie2gmtwszof.onion",
		"-jsonReport",
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out[:]))
	fmt.Println("done executing")
}

