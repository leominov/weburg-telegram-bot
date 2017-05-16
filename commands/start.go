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
			Usage:  "Address to listen on for web interface and telemetry",
			EnvVar: "WEBURG_BOT_LISTEN_ADDR",
		},
		cli.StringFlag{
			Name:   "metrics-path",
			Usage:  "Path under which to expose metrics",
			EnvVar: "WEBURG_BOT_METRICS_PATH",
		},
		cli.StringFlag{
			Name:   "database, db",
			Usage:  "Path to database file",
			EnvVar: "WEBURG_BOT_DATABASE_PATH",
		},
		cli.BoolFlag{
			Name:   "disable-messenger",
			Usage:  "Disable sending messages",
			EnvVar: "WEBURG_BOT_DISABLE_MESSENGER",
		},
		cli.StringFlag{
			Name:   "config-file",
			Value:  "./config.yaml",
			Usage:  "Configuration file",
			EnvVar: "WEBURG_BOT_CONFIG_FILE",
		},
		DebugFlag,
		NoColorFlag,
	},
	Action: func(c *cli.Context) {
		HandleGlobalFlags(GlobalFlagsConfig{
			Debug:   c.Bool("debug"),
			NoColor: c.Bool("no-color"),
		})

		logrus.Infof("Starting %s %s...", c.App.Name, c.App.Version)

		config := bot.NewConfig()
		if len(c.String("config-file")) != 0 {
			logrus.Infof("Loading configuration from file '%s'...", c.String("config-file"))
			if err := config.LoadFromFile(c.String("config-file")); err != nil {
				logrus.Fatal(err)
			}
		}

		config.LoadFromContext(c)

		logrus.Infof("Messenger disabled: %v", config.DisableMessenger)
		logrus.Debugf("Configuration: %s", config.ToString())

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
