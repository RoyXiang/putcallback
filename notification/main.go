package notification

import (
	"os"
	"time"

	"github.com/gregdel/pushover"
)

var (
	isPushoverEnabled bool
	pushoverApp       *pushover.Pushover
	recipient         *pushover.Recipient
)

func init() {
	appToken := os.Getenv("PUSHOVER_APP_TOKEN")
	userToken := os.Getenv("PUSHOVER_USER_TOKEN")
	if appToken != "" && userToken != "" {
		isPushoverEnabled = true
		pushoverApp = pushover.New(appToken)
		recipient = pushover.NewRecipient(userToken)
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
}
