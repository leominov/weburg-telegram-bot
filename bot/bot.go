package bot

import "github.com/leominov/weburg-telegram-bot/metrics"

type Config struct {
	Token       string `json:"token"`
	Watch       bool   `json:"watch"`
	ListenAddr  string `json:"listen_addr"`
	MetricsPath string `json:"metrics_path"`
}

type Bot struct {
	Config       Config
	w            Watcher
	isConfigured bool
}

func New(c Config) *Bot {
	return &Bot{
		Config:       c,
		isConfigured: false,
	}
}

func (b *Bot) Setup() error {
	metrics.InitMetrics()

	telegram := Telegram{
		Token: b.Config.Token,
	}

	if err := telegram.Authorize(); err != nil {
		return err
	}

	watcher := Watcher{
		Telegram: telegram,
	}

	b.w = watcher

	b.isConfigured = true

	return nil
}

func (b *Bot) Start() {
	go metrics.ServeMetrics(b.Config.ListenAddr, b.Config.MetricsPath)
	b.w.Start()
}
