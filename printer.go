package main

import (
  "log"
  flow "github.com/trustmaster/goflow"
)

type Printer struct {
	flow.Component
	Word <-chan HashResult // result of string comparison
}

func (p *Printer) OnWord(res HashResult) {
	log.Printf("%s", res)
}


