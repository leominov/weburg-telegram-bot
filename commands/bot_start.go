package commands

import (
	"github.com/leominov/weburg-telegram-bot/metrics"
	"github.com/leominov/weburg-telegram-bot/watcher"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

type StartCommandConfig struct {
	Token       string
	Watch       bool
	Debug       bool
	NoColor     bool
	ListenAddr  string
	MetricsPath string
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
		cfg := StartCommandConfig{
			Token:       c.String("token"),
			Watch:       c.Bool("watch"),
			Debug:       c.Bool("debug"),
			NoColor:     c.Bool("no-color"),
			ListenAddr:  c.String("listen-address"),
			MetricsPath: c.String("metrics-path"),
		}

		HandleGlobalFlags(cfg)

		w := watcher.Watcher{watcher.Telegram{
			Token: cfg.Token,
		}}

		if err := w.Telegram.Authorize(); err != nil {
			logrus.Fatalf("%+v", err)
		}

		if cfg.Watch {
			go w.Start()
			metrics.ServeMetrics(cfg.ListenAddr, cfg.MetricsPath)
		}
	},
}

func init() {
	metrics.InitMetrics()
}
