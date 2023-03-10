package main

import (
	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	tb "gopkg.in/telebot.v3"
)

type config struct {
	Messengers messengers     `json:"messengers"`
	Bridges    []bridgeConfig `json:"bridges"`
}

type messengers struct {
	Utopia   utopiago.Config `json:"utopia"`
	Telegram telegramConfig  `json:"telegram"`
}

type bridgeConfig struct {
	UtopiaChannelID string `json:"utopiaChannelID"`
	TelegramChatID  int64  `json:"telegramChatID"`
}

type telegramConfig struct {
	BotToken string `json:"botToken"`
}

type bridges struct {
	UtopiaToTelegram map[string]int64
	TelegramToUtopia map[int64]string
}

type bot struct {
	ChatBot *uchatbot.ChatBot
	TgBot   *tb.Bot
	Bridges bridges
}
