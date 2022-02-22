package main

import (
	"github.com/urfave/cli/v2"
	"kt/cmd"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		cmd.List,
		cmd.Shutdown,
		cmd.Start,
	}
	app.Run(os.Args)
}
