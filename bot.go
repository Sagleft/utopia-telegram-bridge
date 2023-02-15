package main

import (
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uchatbot-engine"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	tb "gopkg.in/telebot.v3"
)

const (
	previewTitle       = "bridge"
	configFilePath     = "config.json"
	previewColor       = "green"
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
		log.Printf("build bridge U %q <-> T %v", r.UtopiaChannelID, r.TelegramChatID)
		b.Bridges.UtopiaToTelegram[r.UtopiaChannelID] = r.TelegramChatID
		b.Bridges.TelegramToUtopia[r.TelegramChatID] = r.UtopiaChannelID
	}

	return b
}

func (b *bot) setChatBot(cb *uchatbot.ChatBot) {
	b.ChatBot = cb
}

func (b *bot) setTelegramBot(tgBot *tb.Bot) {
	b.TgBot = tgBot
}

func main() {
	figure.NewColorFigure(previewTitle, "", previewColor, true).Print()
	fmt.Println()

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
	cb, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.Messengers.Utopia,
		Chats:  chats,
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        onContactMessage,
			OnChannelMessage:        b.onUtopiaChannelMessage,
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

	b.setChatBot(cb)
	b.setTelegramBot(tgBot)

	tgBot.Handle(tb.OnText, b.onTelegramMessage)

	go tgBot.Start()

	log.Println("bot started")
	swissknife.RunInBackground()
}

func (b *bot) getTelegramBridge(chatID int64) (string, bool) {
	br, isExists := b.Bridges.TelegramToUtopia[chatID]
	return br, isExists
}

func (b *bot) getUtopiaBridge(channelID string) (int64, bool) {
	br, isExists := b.Bridges.UtopiaToTelegram[channelID]
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

func (b *bot) onUtopiaChannelMessage(m structs.WsChannelMessage) {
	chatID, isExists := b.getUtopiaBridge(m.ChannelID)
	if !isExists {
		log.Printf("unknown utopia channel ID %v, bridge not found", chatID)
		return
	}

	if err := b.sendToTelegram(chatID, m.Nick, m.Text); err != nil {
		color.Red("send message to telegram: %s", err.Error())
	}
}

func (b *bot) sendToTelegram(chatID int64, nickname string, message string) error {
	_, err := b.TgBot.Send(
		tb.ChatID(chatID),
		fmt.Sprintf("%s: %s", nickname, message),
	)
	return err
}

func (b *bot) sendToUtopia(channelID string, nickname string, message string) error {
	return b.ChatBot.SendChannelMessage(
		channelID,
		fmt.Sprintf("%s: %s", nickname, message),
	)
}
