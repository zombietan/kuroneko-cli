package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/zombietan/kuroneko/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "kuroneko"
	app.Usage = "Display delivery status"
	app.Version = "0.0.0"
	app.Action = cmd.Kuroneko
	app.Run(os.Args)
}
