package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/zombietan/kuroneko/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "kuroneko"
	app.Usage = "Display delivery status"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		cli.UintFlag{
			Name:  "serial, s",
			Usage: "連番は10件まで",
		},
	}
	app.Action = func(c *cli.Context) error {
		set := c.GlobalIsSet("serial")
		if set {
			items := c.Uint("serial")
			if 1 <= items && items <= 10 {
				cmd.TrackSerialNumbers(c)
			} else {
				fmt.Println("連番で取得できるのは10件までです")
			}
		} else {
			cmd.TrackNumber(c)
		}

		return nil
	}

	app.Run(os.Args)
}
