package info

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/zijoziks/dzemo/internal/logger"
)

var Command = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Displays basic information about the bot",
}

func Handle(event *events.ApplicationCommandInteractionCreate) {
	err := event.CreateMessage(discord.NewMessageCreate().
		WithContent("Džemo is a utility bot.").
		WithEphemeral(true),
	)
	if err != nil {
		logger.Logger.Error("failed to create a message", "error", err)
		// log.Println("Failed to create a message: ", err)
	}
}
