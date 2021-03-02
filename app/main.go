package main

import (
	"github.com/atom-eight/tmt-backend/app/cmd"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
