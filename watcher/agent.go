package watcher

import (
	"gotel/bot"
	"time"

	"github.com/tucnak/telebot"
)

type RssAgent struct {
	Type           string
	Endpoint       string
	Interval       time.Duration
	CetegoryFilter []string
	Sender         bot.WeburgBot
	Channel        telebot.Chat

	first    bool
	lastGUID string
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
		},
		RssAgent{
			Type:     "music",
			Endpoint: "http://rss.weburg.net/music/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_music",
			},
		},
		RssAgent{
			Type:     "news",
			Endpoint: "http://rss.weburg.net/news/all.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_times",
			},
		},
		RssAgent{
			Type:     "series",
			Endpoint: "http://rss.weburg.net/movies/series.rss",
			Interval: time.Minute,
			Channel: telebot.Chat{
				Type:     "channel",
				Username: "weburg_series",
			},
		},
	}
)
