package main

import (
	"gotel/bot"
	"gotel/clicommand"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "weburg-telegram-bot"
	app.Version = bot.Version()
	app.Author = "Lev Aminov <mailto@levaminov.ru>"

	app.Commands = []cli.Command{
		clicommand.BotStartCommand,
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
