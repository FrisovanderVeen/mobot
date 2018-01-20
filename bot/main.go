package main

import (
	"github.com/FrisovanderVeen/mobot/bot/cmd"
)

func main() {
	app := cmd.NewApp()
	app.RunAndExitOnError()
}
