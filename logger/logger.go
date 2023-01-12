package logger

import (
	"fmt"
	"log"
)

var Loglevel string

type Logstruct struct {
	Loglevel string
}

func (l Logstruct) Logger(title string, t interface{}) {
	if l.Loglevel == "DEBUG" {
		log.Println(title)
		log.Printf("LOG : %v", t)
		fmt.Println("---------------------------")
	}
}
