package main

import (
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uchatbot-engine"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
	tb "gopkg.in/telebot.v3"
)

const (
	configFilePath     = "config.json"
	longPollerInterval = 15 * time.Second
)

func newBot(cfg config) *bot {
	b := &bot{
		Redirects: redirector{
			UtopiaToTelegram: make(map[string]int64),
			TelegramToUtopia: make(map[int64]string),
		},
	}

	for _, r := range cfg.Redirects {
		b.Redirects.UtopiaToTelegram[r.UtopiaChannelID] = r.TelegramChatID
		b.Redirects.TelegramToUtopia[r.TelegramChatID] = r.UtopiaChannelID
	}

	return b
}

func main() {
	cfg := config{}
	if err := swissknife.ParseStructFromJSONFile(configFilePath, &cfg); err != nil {
		color.Red("read config: %s", err.Error())
		return
	}

	b := newBot(cfg)

	// setup utopia bot
	chats := []uchatbot.Chat{}
	for _, r := range cfg.Redirects {
		chats = append(chats, uchatbot.Chat{ID: r.UtopiaChannelID})
	}

	_, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.Utopia,
		Chats:  chats,
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        onContactMessage,
			OnChannelMessage:        b.onChannelMessage,
			OnPrivateChannelMessage: onPrivateChannelMessage,
		},
		UseErrorCallback: true,
		ErrorCallback:    onError,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// setup telegram bot
	tgBot, err := tb.NewBot(tb.Settings{
		Token:  cfg.Telegram.BotToken,
		Poller: getTgPoller(),
	})
	if err != nil {
		log.Fatalf("create tg bot: %v", err)
	}
	go tgBot.Start()

	swissknife.RunInBackground()
}

func (b *bot) onChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[CHANNEL] %s: %s\n", m.Nick, m.Text)
}

func (b *bot) sendToTelegram(chatID int64, nickname string, message string) error {
	// TODO
	return nil
}

func (b *bot) sendToUtopia(channelID string, nickname string, message string) error {
	// TODO
	return nil
}
