package main

import (
	"github.com/ChrisEineke/concourse"
)

func main() {
	app := &concourse.NodeOrchestratorApp{}
	app.Configure()
	app.Run()
}
