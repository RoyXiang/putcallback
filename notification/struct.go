package notification

type notifier interface {
	send(message string)
}
