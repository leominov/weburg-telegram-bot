package clicommand

import (
	"gotel/bot"
	"gotel/watcher"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

type BotStartConfig struct {
	Token    string
	RssWatch bool
	Debug    bool
	NoColor  bool
}

var BotStartCommand = cli.Command{
	Name:  "start",
	Usage: "Starts a Weburg bot",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "token, t",
			Value:  "",
			Usage:  "Your Telegram API token",
			EnvVar: "WEBURG_BOT_TOKEN",
		},
		cli.BoolFlag{
			Name:   "rss-watch, r",
			Usage:  "Enable RSS watching",
			EnvVar: "WEBURG_BOT_RSS_WATCH",
		},
		DebugFlag,
		NoColorFlag,
	},
	Action: func(c *cli.Context) {
		cfg := BotStartConfig{
			Token:    c.String("token"),
			RssWatch: c.Bool("rss-watch"),
			Debug:    c.Bool("debug"),
			NoColor:  c.Bool("no-color"),
		}

		HandleGlobalFlags(cfg)

		w := bot.WeburgBot{
			Token: cfg.Token,
		}

		if err := w.Start(); err != nil {
			logrus.Fatalf("%+v", err)
		}

		if cfg.RssWatch {
			w := watcher.RssWatcher{
				Sender: w,
			}

			go w.StartWatch()
		}

		w.Listen()
	},
}
