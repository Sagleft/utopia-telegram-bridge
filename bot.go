package main

import (
	"fmt"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
	tb "gopkg.in/telebot.v3"
)

const (
	configFilePath     = "config.json"
	longPollerInterval = 15 * time.Second
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

func main() {
	cfg := config{}
	if err := swissknife.ParseStructFromJSONFile(configFilePath, &cfg); err != nil {
		color.Red("read config: %s", err.Error())
		return
	}

	// setup utopia bot
	_, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: cfg.Utopia,
		Chats:  cfg.Chats,
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        OnContactMessage,
			OnChannelMessage:        OnChannelMessage,
			OnPrivateChannelMessage: OnPrivateChannelMessage,
		},
		UseErrorCallback: true,
		ErrorCallback:    onError,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// setup telegram bot
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Telegram.BotToken,
		Poller: getTgPoller(),
	})
	if err != nil {
		log.Fatalf("create tg bot: %v", err)
	}
	go b.Start()

	swissknife.RunInBackground()
}

func OnContactMessage(m structs.InstantMessage) {
	fmt.Printf("[CONTACT] %s: %s\n", m.Nick, m.Text)
}

func OnChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[CHANNEL] %s: %s\n", m.Nick, m.Text)
}

func OnPrivateChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[PRIVATE] [%s] %s: %s\n", m.ChannelName, m.Nick, m.Text)
}

func onError(err error) {
	color.Red(err.Error())
}
