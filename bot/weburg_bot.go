package bot

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"
)

type WeburgBot struct {
	b *telebot.Bot

	Token     string
	Channel   string
	StartTime time.Time
}

func (w *WeburgBot) Authorize() error {
	wbot, err := telebot.NewBot(w.Token)
	if err != nil {
		return err
	}

	w.StartTime = time.Now()

	logrus.Info("Authorized as ", wbot.Identity.Username)

	w.b = wbot

	return nil
}

func (w *WeburgBot) Send(c telebot.Chat, message string) error {
	logrus.WithField("channel", c.Username).Debug(message)
	return w.b.SendMessage(c, message, nil)
}
