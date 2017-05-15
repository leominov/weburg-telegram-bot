package bot

import (
	"errors"
	"sync"

	"github.com/leominov/weburg-telegram-bot/metrics"
)

type Config struct {
	Token       string `json:"token"`
	Watch       bool   `json:"watch"`
	ListenAddr  string `json:"listen_addr"`
	MetricsPath string `json:"metrics_path"`
}

type Bot struct {
	Config       Config
	m            *Messenger
	isConfigured bool
	stopChan     chan bool
	doneChan     chan bool
}

func New(c Config) *Bot {
	return &Bot{
		Config:       c,
		isConfigured: false,
		stopChan:     make(chan bool),
		doneChan:     make(chan bool),
	}
}

func (b *Bot) Setup() error {
	metrics.InitMetrics()

	messenger := &Messenger{
		Token: b.Config.Token,
	}

	if err := messenger.Authorize(); err != nil {
		return err
	}

	b.m = messenger
	b.isConfigured = true

	return nil
}

func (b *Bot) Start() error {
	var wg sync.WaitGroup

	if !b.isConfigured {
		return errors.New("Must be configured before start")
	}

	go metrics.ServeMetrics(b.Config.ListenAddr, b.Config.MetricsPath)

	totalAgents := len(AgentsCollection)
	wg.Add(totalAgents)

	for i := 0; i <= totalAgents-1; i++ {
		go func(i int) {
			AgentsCollection[i].Start(b.m)
			wg.Done()
		}(i)
	}

	wg.Wait()

	// Waiting for complete stop
	<-b.doneChan

	return nil
}

func (b *Bot) Stop() error {
	for _, agent := range AgentsCollection {
		agent.Stop()
	}
	close(b.doneChan)
	return nil
}
