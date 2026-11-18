package main

import (
	"github.com/ChrisEineke/concourse"
)

func main() {
	app := &concourse.DeployerApp{}
	app.Configure()
	defer app.Stop()
	app.Run()
}
