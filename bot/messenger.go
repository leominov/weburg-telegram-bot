package bot

import (
	"github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"
)

type Messenger struct {
	b *telebot.Bot

	Token string
}

func (m *Messenger) Authorize() error {
	wbot, err := telebot.NewBot(m.Token)
	if err != nil {
		return err
	}

	logrus.Info("Authorized as ", wbot.Identity.Username)

	m.b = wbot

	return nil
}

func (m *Messenger) Send(c telebot.Chat, message string) error {
	logrus.WithField("channel", c.Username).Debug(message)
	return m.b.SendMessage(c, message, nil)
}
