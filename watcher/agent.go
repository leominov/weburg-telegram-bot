package watcher

import (
	"time"

	"github.com/leominov/weburg-telegram-bot/bot"

	"github.com/tucnak/telebot"
)

type RssAgent struct {
	Type           string
	Endpoint       string
	Interval       time.Duration
	CetegoryFilter []string
	Sender         bot.WeburgBot
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
