package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/zijoziks/dzemo/internal/logger"
)

func main() {
	if err := run(); err != nil {
		logger.Logger.Error("startup failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := loadConfig()
	if err != nil {
		// TODO I was thinking of only outputting WARN if DZEMO_GUILD_ID is not set
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := newClient(config)
	if err != nil {
		return fmt.Errorf("failed to create disgo client: %w", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client.Close(ctx)
	}()

	if _, err = client.Rest.SetGuildCommands(client.ApplicationID, config.guildID, commands); err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.OpenGateway(ctx); err != nil {

		return fmt.Errorf("failed to open disgo gateway: %w", err)
	}

	waitForShutdown()
	return nil
}
