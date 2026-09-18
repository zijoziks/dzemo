// TODO HAVE A LOGGING FILE/SYSTEM WHERE WE CAN SEE WHO VERIFIED WHO WITH A COMMAND IN A DISCORD CHANNEL OR WITH A LOOKUP COMMAND

package verify

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/zijoziks/dzemo/internal/logger"
)

var Command = discord.SlashCommandCreate{
	Name:                     "verify",
	Description:              "Verifikuje korisnika i daje mu pristup serveru",
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionManageRoles),
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "korisnik",
			Description: "Tag korisnika kojeg želite pustiti",
			Required:    true,
		},
	},
}

func NewHandler(memberRoleID snowflake.ID) func(event *events.ApplicationCommandInteractionCreate) {
	return func(event *events.ApplicationCommandInteractionCreate) {
		guildID := event.GuildID()
		if guildID == nil {
			logger.Logger.Error("verify command used outside of a guild")
			return
		}

		target := event.SlashCommandInteractionData().User("korisnik")

		if err := event.Client().Rest.AddMemberRole(*guildID, target.ID, memberRoleID); err != nil {
			logger.Logger.Error("failed to add member role", "error", err, "target", target.Username)
			if err := event.CreateMessage(discord.NewMessageCreate().
				WithContent("Nisam uspio verifikovat korisnika " + target.Mention() + ".").
				WithEphemeral(true),
			); err != nil {
				logger.Logger.Error("failed to create a message", "error", err)
			}
			return
		}

		if err := event.CreateMessage(discord.NewMessageCreate().
			WithContent("Korisnik " + target.Mention() + " je uspješno verifikovan.").
			WithEphemeral(true),
		); err != nil {
			logger.Logger.Error("failed to create a message", "error", err)
		}
	}
}
