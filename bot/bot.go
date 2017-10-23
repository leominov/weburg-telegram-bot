package bot

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

var StateBucket = []byte("statev1")

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
	b.InitMetrics()

	if len(b.Config.Agents) == 0 {
		return errors.New("Agents list cant be empty")
	}

	messenger := &Messenger{
		Token:    b.Config.Token,
		Disabled: b.Config.DisableMessenger,
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

	go b.ServeMetrics()

	totalAgents := len(b.Config.Agents)
	wg.Add(totalAgents)

	for i := 0; i <= totalAgents-1; i++ {
		go func(i int) {
			state, err := b.RestoreStateFor(b.Config.Agents[i].Name)
			if err != nil {
				state = []string{}
			}
			if err := b.Config.Agents[i].Start(b.m, state); err != nil {
				logrus.Error(err)
			}
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
	for _, agent := range b.Config.Agents {
		agent.Stop()
		if err := b.SaveStateFor(agent.Name, agent.lastGuids); err != nil {
			logrus.Error(err)
		}
	}
	if err := b.DB.Close(); err != nil {
		return err
	}
	close(b.doneChan)
	return nil
}
