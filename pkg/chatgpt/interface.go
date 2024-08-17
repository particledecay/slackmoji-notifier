package chatgpt

// ClientInterface defines the methods that a ChatGPT client should implement
type ClientInterface interface {
	GenerateCompletion(prompt string, streamToStdout bool) (string, error)
}
