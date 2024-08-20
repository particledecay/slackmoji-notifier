<p align="center">
  <h1 align="center">Slackmoji Notifier</h1>
  <p align="center"><img src="assets/logo.png" width="128px" height="128px"></p>
  <p align="center">A fun and interactive Slack bot that notifies users about new emoji additions with AI-generated descriptions.</p>
</p>

## Description

Slackmoji Notifier is a Slack bot that monitors your workspace for new emoji additions. When a new emoji is added, it uses OpenAI's GPT model to generate a fun, creative description and sends a notification to a specified Slack channel.

## Features

- Real-time monitoring of new emoji additions in your Slack workspace
- AI-generated descriptions for each new emoji using OpenAI's GPT model
- Customizable Slack channel for notifications
- Easy deployment using Helm charts for Kubernetes

## Installation

### Prerequisites

- Kubernetes cluster
- Helm 3+
- Slack Bot Token and App Token
- OpenAI API Key

### Helm Chart Installation

Install the chart:

   ```bash
   helm install slackmoji-notifier ./chart \
     --set slack.botToken="your-slack-bot-token" \
     --set slack.appToken="your-slack-app-token" \
     --set slack.channel="#your-notification-channel" \
     --set openai.apiKey="your-openai-api-key" \
     --set openai.model="your-preferred-gpt-model" \
     --set verbose=true
   ```

### Run it locally

Clone the repository and install dependencies

    ```bash
    git clone https://github.com/particledecay/slackmoji-notifier
    go mod download
    ```

Build and run the application

    ```bash
    go build -o slackmoji-notifier .
    ./slackmoji-notifier
    ```

## Configuration

Key configuration options:

- Helm values
    - `slack.channel`: The Slack channel where notifications will be sent
    - `slack.botToken`: Your Slack Bot Token
    - `slack.appToken`: Your Slack App Token
    - `openai.apiKey`: Your OpenAI API Key
    - `verbose`: Enable verbose logging (default: false)
- Environment variables
    - `SLACK_CHANNEL`: The Slack channel where notifications will be sent
    - `SLACK_BOT_TOKEN`: Your Slack Bot Token
    - `SLACK_APP_TOKEN`: Your Slack App Token
    - `OPENAI_API_KEY`: Your OpenAI API Key

For more configuration options, see the [values.yaml](./values.yaml) file.

## Add a custom Slack bot to your workspace

1. Create a new Slack app at [api.slack.com/apps](https://api.slack.com/apps) and click "Create New App"
2. Choose "From scratch"
3. Give it a good name and select your workspace
4. Scroll down and give it the icon at [assets/logo.png](./assets/logo.png)
5. Give it the background color '#6c5994'
6. Click "Save Changes"
7. Click "Socket Mode" in the left sidebar
8. Click "Enable Socket Mode" and click "Generate" in the popup (this is your `SLACK_APP_TOKEN`)
9. Click "OAuth & Permissions" in the left sidebar
10. Under "Bot Token Scopes" click "Add an OAuth Scope" and give it the following:
    - `channels:read`
    - `chat:write`
    - `chat:write.public`
    - `emoji:read`
11. Under "OAuth Tokens" click "Install to <Workspace>" and click "Allow"
12. Copy the "Bot User OAuth Token" (this is your `SLACK_BOT_TOKEN`)
13. Run the application locally (or within a Kubernetes cluster) and set `SLACK_CHANNEL` to any public channel

## Why?

Emojis are a fun and expressive part of Slack communication. Slackmoji Notifier adds an extra layer of enjoyment by:

- Ensuring no new emoji goes unnoticed
- Providing funny, sometimes nonsensical AI-generated sentences
- Encouraging emoji usage and creativity within your team

## Known Issues

Check out the [Issues](https://github.com/particledecay/slackmoji-notifier/issues) section or specifically [issues created by me](https://github.com/particledecay/slackmoji-notifier/issues?q=is:issue+is:open+sort:updated-desc+author:particledecay)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
