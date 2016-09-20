package watcher

import (
	"time"

	"gotel/bot"

	"github.com/Sirupsen/logrus"
	rss "github.com/ungerik/go-rss"
)

func (r *RssAgent) CanPost(item rss.Item) bool {
	if r.lastGUID == item.GUID {
		return false
	} else if len(r.CetegoryFilter) == 0 {
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

func (r *RssAgent) Start(sender bot.WeburgBot) error {
	r.first = true
	r.Sender = sender

	feed, err := rss.Read(r.Endpoint)
	if err != nil {
		return err
	}

	logrus.Infof("Found feed '%s'", feed.Title)

	if len(feed.Item) > 0 {
		r.lastGUID = feed.Item[0].GUID
	}

	for {
		feed, err = rss.Read(r.Endpoint)
		if err != nil {
			logrus.Errorf("Error with %s: %+v", r.Endpoint, err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = r.itemHandler(feed.Item)
		if err != nil {
			logrus.Errorf("Error with %s: %+v", r.Endpoint, err)
		}

		<-time.After(r.Interval)
	}

	return nil
}

func (r *RssAgent) itemHandler(items []rss.Item) error {
	logrus.Debugf("Got %d items in '%s' channel", len(items), r.Type)

	if len(items) == 0 || r.first == true {
		logrus.Debugf("Skipping update in '%s' channel", r.Type)
		r.first = false
		return nil
	}

	item := items[0]
	if r.CanPost(item) != true {
		return nil
	}

	r.lastGUID = item.GUID

	logrus.Infof("Send '%s' to %s channel", item.Title, r.Type)
	return r.Sender.SendMessage(r.Channel, item.Title+"\n\n"+item.Link)
}
