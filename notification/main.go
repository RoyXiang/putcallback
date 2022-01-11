package notification

import (
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gregdel/pushover"
)

var (
	isPushoverEnabled bool
	pushoverApp       *pushover.Pushover
	recipient         *pushover.Recipient

	tgBotApi *tgbotapi.BotAPI
	tgChatId int64
)

func init() {
	appToken := os.Getenv("PUSHOVER_APP_TOKEN")
	userToken := os.Getenv("PUSHOVER_USER_TOKEN")
	if appToken != "" && userToken != "" {
		isPushoverEnabled = true
		pushoverApp = pushover.New(appToken)
		recipient = pushover.NewRecipient(userToken)
	}

	var err error
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatId := os.Getenv("TELEGRAM_CHAT_ID")
	if botToken != "" && telegramChatId != "" {
		tgChatId, err = strconv.ParseInt(telegramChatId, 10, 64)
		if err == nil {
			tgBotApi, _ = tgbotapi.NewBotAPI(botToken)
		}
	}
}

func Send(message string) {
	if isPushoverEnabled {
		msg := &pushover.Message{
			Title:     "Put.io",
			Message:   message,
			Sound:     pushover.SoundNone,
			Timestamp: time.Now().Unix(),
		}
		_, _ = pushoverApp.SendMessage(msg, recipient)
	}
	if tgBotApi != nil {
		msg := tgbotapi.NewMessage(tgChatId, message)
		_, _ = tgBotApi.Send(msg)
	}
}
