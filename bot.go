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
		Bridges: bridges{
			UtopiaToTelegram: make(map[string]int64),
			TelegramToUtopia: make(map[int64]string),
		},
	}

	for _, r := range cfg.Bridges {
		b.Bridges.UtopiaToTelegram[r.UtopiaChannelID] = r.TelegramChatID
		b.Bridges.TelegramToUtopia[r.TelegramChatID] = r.UtopiaChannelID
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
	for _, r := range cfg.Bridges {
		chats = append(chats, uchatbot.Chat{ID: r.UtopiaChannelID})
	}

	_, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.Messengers.Utopia,
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
		Token:  cfg.Messengers.Telegram.BotToken,
		Poller: getTgPoller(),
	})
	if err != nil {
		log.Fatalf("create tg bot: %v", err)
	}

	tgBot.Handle(tb.OnText, b.onTelegramMessage)

	go tgBot.Start()
	swissknife.RunInBackground()
}

func (b *bot) getTelegramBridge(chatID int64) (string, bool) {
	br, isExists := b.Bridges.TelegramToUtopia[chatID]
	return br, isExists
}

func (b *bot) onTelegramMessage(c tb.Context) error {
	var (
		user   = c.Sender()
		text   = c.Text()
		chatID = c.Chat().ID
	)

	uChannelID, isExists := b.getTelegramBridge(chatID)
	if !isExists {
		log.Printf("unknown telegram chat ID %v, bridge not found", chatID)
		return nil
	}

	nickname := getTelegramNickname(user)
	if err := b.sendToUtopia(uChannelID, nickname, text); err != nil {
		return fmt.Errorf("send message to utopia: %w", err)
	}

	return nil
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
