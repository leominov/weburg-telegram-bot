package bot

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/leominov/weburg-telegram-bot/metrics"
)

var StateBucket = []byte("statev1")

type Config struct {
	Token        string `json:"token"`
	Watch        bool   `json:"watch"`
	ListenAddr   string `json:"listen_addr"`
	MetricsPath  string `json:"metrics_path"`
	DatabasePath string `json:"database_path"`
}

type Bot struct {
	Config       *Config
	DB           *bolt.DB
	m            *Messenger
	isConfigured bool
	stopChan     chan bool
	doneChan     chan bool
}

func New(c *Config) *Bot {
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

	b.m = messenger

	if err := messenger.Authorize(); err != nil {
		return err
	}

	db, err := bolt.Open(b.Config.DatabasePath, 0600, nil)
	if err != nil {
		return err
	}

	b.DB = db
	b.isConfigured = true

	return b.DB.Update(func(tx *bolt.Tx) error {
		// Always create State bucket.
		if _, err := tx.CreateBucketIfNotExists(StateBucket); err != nil {
			return err
		}
		return nil
	})
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
			state, err := b.RestoreStateFor(AgentsCollection[i].Type)
			if err != nil {
				state = []string{}
			}
			AgentsCollection[i].Start(b.m, state)
			wg.Done()
		}(i)
	}

	wg.Wait()

	// Waiting for complete stop
	<-b.doneChan

	return nil
}

func (b *Bot) RestoreStateFor(agent string) ([]string, error) {
	var state []string
	return state, b.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(StateBucket)
		v := b.Get([]byte(agent))
		if len(v) == 0 {
			return errors.New("Nothing found")
		}
		if err := json.Unmarshal(v, &state); err != nil {
			return err
		}
		return nil
	})
}

func (b *Bot) SaveStateFor(agent string, state []string) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(StateBucket)
		encoded, err := json.Marshal(state)
		if err != nil {
			return err
		}
		return b.Put([]byte(agent), encoded)
	})
}

func (b *Bot) Stop() error {
	for _, agent := range AgentsCollection {
		agent.Stop()
		if err := b.SaveStateFor(agent.Type, agent.lastGuids); err != nil {
			logrus.Error(err)
		}
	}
	if err := b.DB.Close(); err != nil {
		return err
	}
	close(b.doneChan)
	return nil
}
