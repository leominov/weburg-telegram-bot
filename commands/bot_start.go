package commands

import (
	"github.com/Sirupsen/logrus"
	"github.com/leominov/weburg-telegram-bot/bot"

	"github.com/codegangsta/cli"
)

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
			Name:   "watch, w",
			Usage:  "Enable RSS watching",
			EnvVar: "WEBURG_BOT_WATCH",
		},
		cli.StringFlag{
			Name:   "listen-address",
			Value:  ":9109",
			Usage:  "Address to listen on for web interface and telemetry",
			EnvVar: "WEBURG_BOT_LISTEN_ADDR",
		},
		cli.StringFlag{
			Name:   "metrics-path",
			Value:  "/metrics",
			Usage:  "Path under which to expose metrics",
			EnvVar: "WEBURG_BOT_METRICS_PATH",
		},
		DebugFlag,
		NoColorFlag,
	},
	Action: func(c *cli.Context) {
		HandleGlobalFlags(GlobalFlagsConfig{
			Debug:   c.Bool("debug"),
			NoColor: c.Bool("no-color"),
		})

		config := bot.Config{
			Token:       c.String("token"),
			Watch:       c.Bool("watch"),
			ListenAddr:  c.String("listen-address"),
			MetricsPath: c.String("metrics-path"),
		}

		logrus.Infof("Starting %s %s...", c.App.Name, c.App.Version)

		b := bot.New(config)
		if err := b.Setup(); err != nil {
			logrus.Fatal(err)
		}

		if !b.Config.Watch {
			logrus.Printf("Configuration: %+v", b.Config)
			return
		}

		if err := b.Start(); err != nil {
			logrus.Fatal(err)
		}
	},
}
