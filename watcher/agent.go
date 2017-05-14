package watcher

import (
	"errors"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/leominov/weburg-telegram-bot/metrics"

	"github.com/tucnak/telebot"
	rss "github.com/ungerik/go-rss"
)

const (
	DefaultCacheSize = 1
)

type Agent struct {
	Type           string
	Endpoint       string
	Interval       time.Duration
	CetegoryFilter []string
	Telegram       Telegram
	Channel        telebot.Chat
	CacheSize      int

	firstPoll bool
	lastGuids []string
}

var (
	AgentsCollection = []Agent{
		Agent{
			Type:     "movies",
			Endpoint: "http://rss.weburg.net/movies/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_movies",
			},
			CacheSize: 3,
		},
		Agent{
			Type:     "music",
			Endpoint: "http://rss.weburg.net/music/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_music",
			},
			CacheSize: 3,
		},
		Agent{
			Type:     "news",
			Endpoint: "http://rss.weburg.net/news/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_times",
			},
			CacheSize: 10,
		},
		Agent{
			Type:     "series",
			Endpoint: "http://rss.weburg.net/movies/series.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_series",
			},
			CacheSize: 2,
		},
	}
)

func (a *Agent) CanPost(item rss.Item) bool {
	for _, guid := range a.lastGuids {
		if item.GUID == guid {
			return false
		}
	}

	if len(a.CetegoryFilter) == 0 {
		return true
	}

	for _, filterCategory := range a.CetegoryFilter {
		for _, category := range item.Category {
			if filterCategory == category {
				return true
			}
		}
	}

	return false
}

func (a *Agent) CacheItems(items []rss.Item) error {
	if len(items) == 0 {
		return errors.New("Empty items list")
	}

	a.lastGuids = []string{}
	for _, item := range items {
		if len(a.lastGuids) == a.CacheSize {
			break
		}
		a.lastGuids = append(a.lastGuids, item.GUID)
	}

	logrus.Debugf("Update cached '%s' GUIDs list (max.: %d): %s", a.Type, a.CacheSize, strings.Join(a.lastGuids, ", "))

	return nil
}

func (a *Agent) Start(telegram Telegram) error {
	a.firstPoll = true
	a.Telegram = telegram
	a.lastGuids = []string{}

	if a.CacheSize == 0 {
		a.CacheSize = DefaultCacheSize
	}

	metrics.PullsTotalCounter.Inc()
	metrics.PullsTotalCounters[a.Type].Inc()
	feed, err := rss.Read(a.Endpoint)
	if err != nil {
		metrics.PullsFailCounter.Inc()
		metrics.PullsFailCounters[a.Type].Inc()
		return err
	}

	logrus.Infof("Found feed '%s'", feed.Title)

	a.CacheItems(feed.Item)

	for {
		metrics.PullsTotalCounter.Inc()
		metrics.PullsTotalCounters[a.Type].Inc()

		feed, err = rss.Read(a.Endpoint)
		if err != nil {
			metrics.PullsFailCounter.Inc()
			metrics.PullsFailCounters[a.Type].Inc()
			logrus.Errorf("Error with %s: %+v", a.Endpoint, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := a.Process(feed.Item); err != nil {
			logrus.Errorf("Error with %s: %+v", a.Endpoint, err)
		}

		<-time.After(a.Interval)
	}

	return nil
}

func (a *Agent) Process(items []rss.Item) error {
	var checks int
	var changed bool

	logrus.Debugf("Got %d items in '%s' channel", len(items), a.Type)

	if len(items) == 0 || a.firstPoll == true {
		logrus.Debugf("Skipping update in '%s' channel", a.Type)
		a.firstPoll = false
		return nil
	}

	for _, item := range items {
		if checks == a.CacheSize {
			break
		}
		if a.CanPost(item) == true {
			changed = true
			if err := a.Notify(item); err != nil {
				logrus.Error(err)
			}
		}
		checks++
	}

	if changed {
		a.CacheItems(items)
	}

	return nil
}

func (a *Agent) Notify(item rss.Item) error {
	logrus.Infof("Send '%s' to %s channel", item.Title, a.Type)

	metrics.MessagesTotalCounter.Inc()
	metrics.MessagesTotalCounters[a.Type].Inc()

	if err := a.Telegram.Send(a.Channel, item.Title+"\n\n"+item.Link); err != nil {
		metrics.MessagesFailCounter.Inc()
		metrics.MessagesFailCounters[a.Type].Inc()

		return err
	}

	return nil
}
