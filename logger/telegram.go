package logger

import (
	"fmt"
	"log"
	"os"

	tgram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramLogger struct {
	bot    *tgram.BotAPI
	chatID int64
}

func NewTelegram(token string, chatID int64) (*TelegramLogger, error) {
	bot, err := tgram.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	logger := TelegramLogger{bot, chatID}
	log.SetOutput(&logger)
	return &logger, nil
}

func (t *TelegramLogger) Write(data []byte) (int, error) {
	str := string(data)
	m := tgram.NewMessage(t.chatID, str)
	_, err := t.bot.Send(m)
	n, stdErr := os.Stderr.WriteString(str)
	return n, fmt.Errorf("%v: %v", err, stdErr)
}
