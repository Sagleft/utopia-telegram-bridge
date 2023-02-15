package main

import (
	"fmt"
	"log"

	swissknife "github.com/Sagleft/swiss-knife"
	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
)

const (
	configFilePath = "config.json"
)

type config struct {
	Utopia utopiago.Config `json:"utopia"`
	Chats  []uchatbot.Chat `json:"chats"`
}

func main() {
	cfg := config{}
	if err := swissknife.ParseStructFromJSONFile(configFilePath, &cfg); err != nil {
		color.Red("read config: %s", err.Error())
		return
	}

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
