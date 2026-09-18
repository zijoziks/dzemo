package welcome

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/zijoziks/dzemo/internal/logger"
)

// TODO ADD CONFIGURABILITY

func NewListener(channelID snowflake.ID) func(event *events.GuildMemberJoin) {
	return func(event *events.GuildMemberJoin) {
		_, err := event.Client().Rest.CreateMessage(channelID, discord.NewMessageCreate().
			WithContent(event.Member.Mention()+"\n\n🇧🇦 Dobro dosli u server! "+
				"Recite nesto o sebi kako bi ste dobili pristup ostatku server.\n\n"+
				"🇬🇧 Welcome to the server! "+
				"Tell us a bit about yourself to get access to rest of the server."),
		)
		if err != nil {
			logger.Logger.Error("failed to send welcome message", "error", err)
		} else {
			logger.Logger.Info("welcome message sent")
		}
	}
}
