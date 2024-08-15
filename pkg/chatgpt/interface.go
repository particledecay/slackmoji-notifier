package chatgpt

type ClientInterface interface {
	SendMessage(prompt string) (string, error)
}
