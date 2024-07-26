package notification

import (
	"errors"
	"time"

	"github.com/gregdel/pushover"
)

type pushoverNotifier struct {
	pushoverApp *pushover.Pushover
	recipient   *pushover.Recipient
}

func newPushoverNotifier(appToken, userToken string) (*pushoverNotifier, error) {
	if appToken == "" || userToken == "" {
		return nil, errors.New("the tokens of Pushover app should not be empty")
	}
	return &pushoverNotifier{
		pushoverApp: pushover.New(appToken),
		recipient:   pushover.NewRecipient(userToken),
	}, nil
}

func (n *pushoverNotifier) send(message string) {
	msg := &pushover.Message{
		Title:     "Put.io",
		Message:   message,
		Sound:     pushover.SoundNone,
		Timestamp: time.Now().Unix(),
	}
	_, _ = n.pushoverApp.SendMessage(msg, n.recipient)
}
