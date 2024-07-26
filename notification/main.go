package notification

import (
	"log"
	"os"
)

var notifiers []notifier

func init() {
	appToken := os.Getenv("PUSHOVER_APP_TOKEN")
	userToken := os.Getenv("PUSHOVER_USER_TOKEN")
	if n, err := newPushoverNotifier(appToken, userToken); err == nil {
		notifiers = append(notifiers, n)
		log.Print("Pushover notification enabled")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatId := os.Getenv("TELEGRAM_CHAT_ID")
	if n, err := newTelegramNotifier(botToken, telegramChatId); err == nil {
		notifiers = append(notifiers, n)
		log.Print("Telegram bot notification enabled")
	}
}

func Send(message string) {
	log.Print(message)
	for _, n := range notifiers {
		n.send(message)
	}
}
