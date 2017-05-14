package bot

import (
	"github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"
)

type Telegram struct {
	b *telebot.Bot

	Token string
}

func (t *Telegram) Authorize() error {
	wbot, err := telebot.NewBot(t.Token)
	if err != nil {
		return err
	}

	logrus.Info("Authorized as ", wbot.Identity.Username)

	t.b = wbot

	return nil
}

func (t *Telegram) Send(c telebot.Chat, message string) error {
	logrus.WithField("channel", c.Username).Debug(message)
	return t.b.SendMessage(c, message, nil)
}
