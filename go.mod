module github.com/particledecay/slackmoji-notifier

go 1.24

require (
	github.com/rs/zerolog v1.34.0
	github.com/sashabaranov/go-openai v1.40.5
	github.com/slack-go/slack v0.17.3
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	golang.org/x/sys v0.34.0 // indirect
)

replace github.com/particledecay/slackmoji-notifier => ./
