package bot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/leominov/weburg-telegram-bot/metrics"

	"github.com/tucnak/telebot"
	rss "github.com/ungerik/go-rss"
)

const (
	DefaultCacheSize              = 1
	MessageTemplate               = "%s\n\n%s"
	MessageWithCategoriesTemplate = "%s\n%s\n\n%s"
)

var hashCleaner = strings.NewReplacer(" ", "_", "-", "_")

type Agent struct {
	Type            string
	Endpoint        string
	Interval        time.Duration
	CategoryFilter  []string
	Channel         telebot.Chat
	CacheSize       int
	PrintCategories bool
	messenger       *Messenger
	firstPoll       bool
	lastGuids       []string
	stopChan        chan bool
}

func (a *Agent) CanPost(item rss.Item) bool {
	for _, guid := range a.lastGuids {
		if item.GUID == guid {
			return false
		}
	}

	if len(a.CategoryFilter) == 0 {
		return true
	}

	for _, filterCategory := range a.CategoryFilter {
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

	logrus.Debugf("Update cached '%s' GUID list (max.: %d): %s", a.Type, a.CacheSize, strings.Join(a.lastGuids, ", "))

	return nil
}

func (a *Agent) Start(messenger *Messenger, state []string) error {
	a.messenger = messenger
	a.lastGuids = state
	a.stopChan = make(chan bool)

	if len(a.lastGuids) == 0 {
		a.firstPoll = true
	} else {
		logrus.Debugf("GUID list for '%s' channel loaded from database (max.: %d): %s", a.Type, a.CacheSize, strings.Join(a.lastGuids, ", "))
	}

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

	if len(a.lastGuids) == 0 {
		a.CacheItems(feed.Item)
	}

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

		select {
		case <-a.stopChan:
			break
		case <-time.After(a.Interval):
			continue
		}
	}

	return nil
}

func (a *Agent) Stop() {
	close(a.stopChan)
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
	var message string
	logrus.Infof("Send '%s' to %s channel", item.Title, a.Type)

	metrics.MessagesTotalCounter.Inc()
	metrics.MessagesTotalCounters[a.Type].Inc()

	if a.PrintCategories && len(item.Category) != 0 {
		tmpCat := []string{}
		for _, category := range item.Category {
			tmpCat = append(tmpCat, fmt.Sprintf("#%s", hashCleaner.Replace(category)))
		}
		message = fmt.Sprintf(
			MessageWithCategoriesTemplate,
			item.Title,
			strings.Join(tmpCat, " "),
			item.Link,
		)
	} else {
		message = fmt.Sprintf(
			MessageTemplate,
			item.Title,
			item.Link,
		)
	}

	if err := a.messenger.Send(a.Channel, message); err != nil {
		metrics.MessagesFailCounter.Inc()
		metrics.MessagesFailCounters[a.Type].Inc()

		return err
	}

	return nil
}
