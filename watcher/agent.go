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

type RssAgent struct {
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
	RssAgentsCollection = []RssAgent{
		RssAgent{
			Type:     "movies",
			Endpoint: "http://rss.weburg.net/movies/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_movies",
			},
			CacheSize: 3,
		},
		RssAgent{
			Type:     "music",
			Endpoint: "http://rss.weburg.net/music/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_music",
			},
			CacheSize: 3,
		},
		RssAgent{
			Type:     "news",
			Endpoint: "http://rss.weburg.net/news/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_times",
			},
			CacheSize: 10,
		},
		RssAgent{
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

func (r *RssAgent) CanPost(item rss.Item) bool {
	for _, guid := range r.lastGuids {
		if item.GUID == guid {
			return false
		}
	}

	if len(r.CetegoryFilter) == 0 {
		return true
	}

	for _, filterCategory := range r.CetegoryFilter {
		for _, category := range item.Category {
			if filterCategory == category {
				return true
			}
		}
	}

	return false
}

func (r *RssAgent) CacheItems(items []rss.Item) error {
	if len(items) == 0 {
		return errors.New("Empty items list")
	}

	r.lastGuids = []string{}
	for _, item := range items {
		if len(r.lastGuids) == r.CacheSize {
			break
		}
		r.lastGuids = append(r.lastGuids, item.GUID)
	}

	logrus.Debugf("Update cached '%s' GUIDs list (max.: %d): %s", r.Type, r.CacheSize, strings.Join(r.lastGuids, ", "))

	return nil
}

func (r *RssAgent) Start(telegram Telegram) error {
	r.firstPoll = true
	r.Telegram = telegram
	r.lastGuids = []string{}

	if r.CacheSize == 0 {
		r.CacheSize = DefaultCacheSize
	}

	metrics.PullsTotalCounter.Inc()
	metrics.PullsTotalCounters[r.Type].Inc()
	feed, err := rss.Read(r.Endpoint)
	if err != nil {
		metrics.PullsFailCounter.Inc()
		metrics.PullsFailCounters[r.Type].Inc()
		return err
	}

	logrus.Infof("Found feed '%s'", feed.Title)

	r.CacheItems(feed.Item)

	for {
		metrics.PullsTotalCounter.Inc()
		metrics.PullsTotalCounters[r.Type].Inc()

		feed, err = rss.Read(r.Endpoint)
		if err != nil {
			metrics.PullsFailCounter.Inc()
			metrics.PullsFailCounters[r.Type].Inc()
			logrus.Errorf("Error with %s: %+v", r.Endpoint, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := r.itemHandler(feed.Item); err != nil {
			logrus.Errorf("Error with %s: %+v", r.Endpoint, err)
		}

		<-time.After(r.Interval)
	}

	return nil
}

func (r *RssAgent) itemHandler(items []rss.Item) error {
	var checks int
	var changed bool

	logrus.Debugf("Got %d items in '%s' channel", len(items), r.Type)

	if len(items) == 0 || r.firstPoll == true {
		logrus.Debugf("Skipping update in '%s' channel", r.Type)
		r.firstPoll = false
		return nil
	}

	for _, item := range items {
		if checks == r.CacheSize {
			break
		}
		if r.CanPost(item) == true {
			changed = true
			if err := r.Notify(item); err != nil {
				logrus.Error(err)
			}
		}
		checks++
	}

	if changed {
		r.CacheItems(items)
	}

	return nil
}

func (r *RssAgent) Notify(item rss.Item) error {
	logrus.Infof("Send '%s' to %s channel", item.Title, r.Type)

	metrics.MessagesTotalCounter.Inc()
	metrics.MessagesTotalCounters[r.Type].Inc()

	if err := r.Telegram.Send(r.Channel, item.Title+"\n\n"+item.Link); err != nil {
		metrics.MessagesFailCounter.Inc()
		metrics.MessagesFailCounters[r.Type].Inc()

		return err
	}

	return nil
}
