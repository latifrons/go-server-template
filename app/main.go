package main

import (
	"github.com/latifrons/lbserver/app/cmd"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
