package main

import (
	"fmt"

	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
	tb "gopkg.in/telebot.v3"
)

func tgMessageFilter(upd *tb.Update) bool {
	if upd.Message == nil {
		return true // ignore empty messages
	}

	if upd.Message.Sender.IsBot {
		return false // ignore bots
	}

	if upd.Message.IsService() {
		return false // ignore service messages
	}

	return true
}

func getTgPoller() *tb.MiddlewarePoller {
	poller := &tb.LongPoller{Timeout: longPollerInterval}
	return tb.NewMiddlewarePoller(poller, tgMessageFilter)
}

func onContactMessage(m structs.InstantMessage) {
	fmt.Printf("[CONTACT] %s: %s\n", m.Nick, m.Text)
}

func onPrivateChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[PRIVATE] [%s] %s: %s\n", m.ChannelName, m.Nick, m.Text)
}

func onError(err error) {
	color.Red(err.Error())
}
