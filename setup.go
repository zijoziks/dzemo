package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/zijoziks/dzemo/internal/logger"
	"github.com/zijoziks/dzemo/internal/verify"
	"github.com/zijoziks/dzemo/internal/welcome"
)

// TODO Research to resolve the duplication of the struct for JSON-sake
type Config struct {
	token            string
	guildID          snowflake.ID
	welcomeChannelID snowflake.ID
	memberRoleID     snowflake.ID
}

type appSettings struct {
	GuildID          snowflake.ID `json:"guildID"`
	WelcomeChannelID snowflake.ID `json:"welcomeChannelID"`
	MemberRoleID     snowflake.ID `json:"memberRoleID"`
}

func loadConfig() (Config, error) {
	token, err := loadToken()
	if err != nil {
		return Config{}, err
	}

	settings, err := loadAppSettings()
	if err != nil {
		return Config{}, err
	}

	return Config{
		token:            token,
		guildID:          settings.GuildID,
		welcomeChannelID: settings.WelcomeChannelID,
		memberRoleID:     settings.MemberRoleID,
	}, nil
}

func loadToken() (string, error) {
	if credsDir := os.Getenv("CREDENTIALS_DIRECTORY"); credsDir != "" {
		data, err := os.ReadFile(filepath.Join(credsDir, "dzemo_token"))
		if err != nil {
			return "", fmt.Errorf("failed to read token: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	}

	token := os.Getenv("DZEMO_TOKEN")
	if token == "" {
		return "", errors.New("DZEMO_TOKEN is not set")
	}
	return token, nil
}

func loadAppSettings() (appSettings, error) {
	dir := os.Getenv("STATE_DIRECTORY")
	if dir == "" {
		dir = "."
	}
	path := filepath.Join(dir, "config.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return appSettings{}, fmt.Errorf("error reading config file %s: %w", path, err)
	}

	var settings appSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return appSettings{}, fmt.Errorf("error parsing config file %s: %w", path, err)
	}

	if settings.GuildID == 0 || settings.WelcomeChannelID == 0 || settings.MemberRoleID == 0 {
		return appSettings{}, fmt.Errorf("config file %s is missing one or more required IDs", path)
	}

	logger.Logger.Info(
		"config loaded",
		"path", path,
		"guildID", settings.GuildID,
		"welcomeChannelID", settings.WelcomeChannelID,
		"memberRoleID", settings.MemberRoleID)

	return settings, nil
}

func newClient(cfg Config) (*bot.Client, error) {
	verifyHandler := verify.NewHandler(cfg.memberRoleID)
	welcomeListener := welcome.NewListener(cfg.welcomeChannelID)
	dispatcher := newCommandDispatcher(verifyHandler)

	return disgo.New(cfg.token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuildMembers,
			),
		),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			logger.Logger.Info("bot is connected", "user", e.User.Username)
		}),
		bot.WithEventListenerFunc(dispatcher.Handle),
		bot.WithEventListenerFunc(welcomeListener),
		bot.WithLogger(logger.Logger),
	)
}

func waitForShutdown() {
	logger.Logger.Info("Bot is now running. Press Ctrl-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Logger.Info("Shutting down...")
}
