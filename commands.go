package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/zijoziks/dzemo/internal/verify"

	"github.com/zijoziks/dzemo/internal/info"
	"github.com/zijoziks/dzemo/internal/logger"
)

var commands = []discord.ApplicationCommandCreate{
	info.Command,
	verify.Command,
}

type commandHandler func(event *events.ApplicationCommandInteractionCreate)

type commandDispatcher struct {
	verifyHandler commandHandler
}

func newCommandDispatcher(verifyHandler commandHandler) *commandDispatcher {
	return &commandDispatcher{verifyHandler: verifyHandler}
}

func (d *commandDispatcher) Handle(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	user := event.User()
	logger.Logger.Info("command dispatched", "command", data.CommandName(), "user", user.Username)

	switch data.CommandName() {
	case "info":
		info.Handle(event)
	case "verify":
		d.verifyHandler(event)
	}
}
