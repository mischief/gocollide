package main

import (
	flow "github.com/trustmaster/goflow"
)

type Comparator struct {
	flow.Component

	Word   <-chan HashResult // word to compare
	Result chan<- HashResult // result, if compare is true

	Comp string // string to compare to
}

func NewComparator(comp string) *Comparator {
	c := new(Comparator)

	c.Comp = comp

	return c
}

func (c *Comparator) OnWord(res HashResult) {
	if c.Comp == res.Result[:len(c.Comp)] {
		c.Result <- res
	} else if c.Comp == res.Result[len(res.Result)-len(c.Comp):] {
		c.Result <- res
	}
}
