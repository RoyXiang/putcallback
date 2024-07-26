package notification

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramNotifier struct {
	tgBotApi *tgbotapi.BotAPI
	tgChatId int64
}

func newTelegramNotifier(botToken, chatId string) (*telegramNotifier, error) {
	tgChatId, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		return nil, err
	}
	tgBotApi, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	return &telegramNotifier{
		tgBotApi: tgBotApi,
		tgChatId: tgChatId,
	}, nil
}

func (n *telegramNotifier) send(message string) {
	msg := tgbotapi.NewMessage(n.tgChatId, message)
	_, _ = n.tgBotApi.Send(msg)
}
