package watcher

import (
	"errors"
	"strings"
	"time"

	"github.com/leominov/weburg-telegram-bot/bot"
	"github.com/leominov/weburg-telegram-bot/metrics"

	"github.com/Sirupsen/logrus"
	rss "github.com/ungerik/go-rss"
)

const (
	DefaultCacheSize = 1
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

func (r *RssAgent) Start(sender bot.WeburgBot) error {
	r.firstPoll = true
	r.Sender = sender
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

	if err := r.Sender.Send(r.Channel, item.Title+"\n\n"+item.Link); err != nil {
		metrics.MessagesFailCounter.Inc()
		metrics.MessagesFailCounters[r.Type].Inc()

		return err
	}

	return nil
}
