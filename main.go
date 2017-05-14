package main

import (
	"os"

	"github.com/leominov/weburg-telegram-bot/commands"

	"github.com/codegangsta/cli"
)

var Version string = "1.0-beta"

func main() {
	app := cli.NewApp()
	app.Name = "weburg-telegram-bot"
	app.Version = Version
	app.Author = "Lev Aminov <mailto@levaminov.ru>"

	app.Commands = []cli.Command{
		commands.StartCommand,
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	app.Run(os.Args)
}
