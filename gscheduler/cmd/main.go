package main

import (
	"github.comm/gscheduler"
)

func main() {
	c := gscheduler.New()
	c.Start()
	defer c.Stop()
}
