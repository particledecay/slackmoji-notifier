package slack

// ClientInterface is an interface for the Slack client
type ClientInterface interface {
	ListenForEvents() error
	SendMessage(content MessageContent) error
	Stop()
}
