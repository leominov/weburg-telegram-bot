package watcher

import (
	"sync"
	"time"

	"github.com/tucnak/telebot"
)

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

type Watcher struct {
	Telegram Telegram
}

func (w *Watcher) Start() {
	var wg sync.WaitGroup
	var totalAgents int

	totalAgents = len(AgentsCollection)
	wg.Add(totalAgents)

	for i := 0; i <= totalAgents-1; i++ {
		go func(i int) {
			AgentsCollection[i].Start(w.Telegram)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
