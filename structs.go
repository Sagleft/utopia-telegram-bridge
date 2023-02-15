package main

import (
	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
)

type config struct {
	Utopia   utopiago.Config `json:"utopia"`
	Telegram telegramConfig  `json:"telegram"`

	Chats     []uchatbot.Chat `json:"chats"`
	Redirects []redirect      `json:"redirects"`
}

type redirect struct {
	UtopiaChannelID string `json:"utopiaChannelID"`
	TelegramChatID  int64  `json:"telegramChatID"`
}

type telegramConfig struct {
	BotToken string `json:"botToken"`
}

type redirector struct {
	UtopiaToTelegram map[string]int64
	TelegramToUtopia map[int64]string
}

type bot struct {
	Redirects redirector
}
