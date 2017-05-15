package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/leominov/weburg-telegram-bot/bot"

	"github.com/codegangsta/cli"
)

var StartCommand = cli.Command{
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
		cli.StringFlag{
			Name:   "database, db",
			Value:  "./database.db",
			Usage:  "Path to database file",
			EnvVar: "WEBURG_BOT_DATABASE_PATH",
		},
		cli.BoolFlag{
			Name:   "disable-messenger",
			Usage:  "Disable sending messages",
			EnvVar: "WEBURG_BOT_DISABLE_MESSENGER",
		},
		DebugFlag,
		NoColorFlag,
	},
	Action: func(c *cli.Context) {
		HandleGlobalFlags(GlobalFlagsConfig{
			Debug:   c.Bool("debug"),
			NoColor: c.Bool("no-color"),
		})

		config := &bot.Config{
			Token:            c.String("token"),
			Watch:            c.Bool("watch"),
			ListenAddr:       c.String("listen-address"),
			MetricsPath:      c.String("metrics-path"),
			DatabasePath:     c.String("database"),
			DisableMessenger: c.Bool("disable-messenger"),
		}

		logrus.Infof("Starting %s %s...", c.App.Name, c.App.Version)
		logrus.Infof("Messenger disabled: %v", config.DisableMessenger)

		b := bot.New(config)
		if err := b.Setup(); err != nil {
			logrus.Fatal(err)
		}

		if !b.Config.Watch {
			logrus.Printf("Configuration: %+v", b.Config)
			return
		}

		errChan := make(chan error, 10)

		go func() {
			errChan <- b.Start()
		}()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case err := <-errChan:
				if err != nil {
					logrus.Fatal(err)
				}
			case signal := <-signalChan:
				logrus.Printf("Captured %v. Exiting...", signal)
				if err := b.Stop(); err != nil {
					logrus.Fatal(err)
				}
				logrus.Print("Bye")
				os.Exit(0)
			}
		}
	},
}
