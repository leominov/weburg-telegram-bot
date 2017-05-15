package bot

import (
	"github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"
)

type Messenger struct {
	Token    string
	Disabled bool
	b        *telebot.Bot
}

func (m *Messenger) Authorize() error {
	bot, err := telebot.NewBot(m.Token)
	if err != nil {
		return err
	}

	logrus.Infof("Authorized as %s", bot.Identity.Username)

	m.b = bot

	return nil
}

func (m *Messenger) Send(c telebot.Chat, message string) error {
	logrus.WithField("channel", c.Username).Debug(message)
	if m.Disabled {
		return nil
	}
	return m.b.SendMessage(c, message, nil)
}
