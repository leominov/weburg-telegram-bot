package bot

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/tucnak/telebot"
)

const (
	DefaultCacheSize              = 1
	MessageTemplate               = "%s\n\n%s"
	MessageWithCategoriesTemplate = "%s\n%s\n\n%s"
	HashtagTemplate               = "#%s"
)

var hashCleaner = strings.NewReplacer(" ", "_", "-", "_", "+", "")

type Agent struct {
	Name                   string        `yaml:"name" json:"name"`
	Endpoint               Endpoint      `yaml:"endpoint" json:"endpoint"`
	FilterCategories       []string      `yaml:"filter_categories" json:"filter_categories"`
	SkipCategories         []string      `yaml:"skip_categories" json:"skip_categories"`
	SkipItemWithCategories []string      `yaml:"skip_item_with_categories" json:"skip_item_with_categories"`
	PrintCategories        bool          `yaml:"print_categories" json:"print_categories"`
	PrintDescription       bool          `yaml:"print_description" json:"print_description"`
	Interval               time.Duration `yaml:"interval" json:"interval"`
	Channel                telebot.Chat  `yaml:"channel" json:"channel"`
	CacheSize              int           `yaml:"cache_size" json:"cache_size"`

	messenger *Messenger
	firstPoll bool
	lastGuids []string
	stopChan  chan bool
}

func (a *Agent) ClearCategories(l []string) []string {
	var result []string
	if len(a.SkipCategories) == 0 {
		return l
	}
	for _, b := range l {
		ina := false
		for _, c := range a.SkipCategories {
			if b == c {
				logrus.Infof("%s = %s", b, c)
				ina = true
			}
		}
		if !ina {
			result = append(result, b)
		}
	}
	return result
}

func (a *Agent) FormatCategoryName(name string) string {
	return fmt.Sprintf(HashtagTemplate, hashCleaner.Replace(name))
}

func (a *Agent) CanPost(item EndpointItem) bool {
	for _, guid := range a.lastGuids {
		if item.ID == guid {
			return false
		}
	}

	for _, category := range a.SkipItemWithCategories {
		for _, itemCategory := range item.Categories {
			if category == itemCategory {
				return false
			}
		}
	}

	if len(a.FilterCategories) == 0 {
		return true
	}

	for _, filterCategory := range a.FilterCategories {
		for _, category := range item.Categories {
			if filterCategory == category {
				return true
			}
		}
	}

	return false
}

func (a *Agent) CacheItems(items []EndpointItem) error {
	if len(items) == 0 {
		return errors.New("Empty items list")
	}

	a.lastGuids = []string{}
	for _, item := range items {
		if len(a.lastGuids) == a.CacheSize {
			break
		}
		a.lastGuids = append(a.lastGuids, item.ID)
	}

	logrus.Debugf("Update cached '%s' GUID list (max.: %d): %s", a.Name, a.CacheSize, strings.Join(a.lastGuids, ", "))

	return nil
}

func (a *Agent) Start(messenger *Messenger, state []string) error {
	a.messenger = messenger
	a.lastGuids = state
	a.stopChan = make(chan bool)

	if len(a.lastGuids) == 0 {
		a.firstPoll = true
	} else {
		logrus.Debugf("GUID list for '%s' channel loaded from database (max.: %d): %s", a.Name, a.CacheSize, strings.Join(a.lastGuids, ", "))
	}

	if a.CacheSize == 0 {
		a.CacheSize = DefaultCacheSize
	}

	PullsTotalCounter.Inc()
	PullsTotalCounters[a.Name].Inc()
	itemList, err := a.Endpoint.Read()
	if err != nil {
		PullsFailCounter.Inc()
		PullsFailCounters[a.Name].Inc()
		return err
	}

	logrus.Infof("Found feed '%s'", a.Name)

	if len(a.lastGuids) == 0 {
		a.CacheItems(itemList)
	}

	for {
		PullsTotalCounter.Inc()
		PullsTotalCounters[a.Name].Inc()

		itemList, err := a.Endpoint.Read()
		if err != nil {
			PullsFailCounter.Inc()
			PullsFailCounters[a.Name].Inc()
			logrus.Errorf("Error with %s: %+v", a.Endpoint.URL, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err := a.Process(itemList); err != nil {
			logrus.Errorf("Error with %s: %+v", a.Endpoint.URL, err)
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

func (a *Agent) Process(items []EndpointItem) error {
	var checks int
	var changed bool

	logrus.Debugf("Got %d items in '%s' channel", len(items), a.Name)

	if len(items) == 0 || a.firstPoll == true {
		logrus.Debugf("Skipping update in '%s' channel", a.Name)
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

func (a *Agent) Notify(item EndpointItem) error {
	var message string
	logrus.Infof("Send '%s' to %s channel", item.Title, a.Name)

	MessagesTotalCounter.Inc()
	MessagesTotalCounters[a.Name].Inc()

	if a.PrintDescription && len(item.Description) != 0 {
		item.Title = fmt.Sprintf("%s\n%s", item.Title, item.Description)
	}

	if a.PrintCategories && len(item.Categories) != 0 {
		cleanedCategories := a.ClearCategories(item.Categories)
		tmpCat := []string{}
		for _, category := range cleanedCategories {
			tmpCat = append(tmpCat, a.FormatCategoryName(category))
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
		MessagesFailCounter.Inc()
		MessagesFailCounters[a.Name].Inc()
		return err
	}

	return nil
}
