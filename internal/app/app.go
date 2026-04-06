package app

import (
	"context"
	"fmt"

	"cs-smokes-bot/internal/adapters/storage/postgres"
	tgadapter "cs-smokes-bot/internal/adapters/telegram"
	"cs-smokes-bot/internal/config"
	tgtransport "cs-smokes-bot/internal/transport/telegram"
	"cs-smokes-bot/internal/usecase"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	tg, err := tgadapter.NewClient(cfg.BotToken)
	if err != nil {
		return err
	}

	pool, err := postgres.NewPool(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	sessionRepo := postgres.NewSessionRepository(pool)
	smokeRepo := postgres.NewSmokeRepository(pool)

	navigationService := usecase.NewNavigationService(sessionRepo, smokeRepo, tg)
	handler := tgtransport.NewHandler(tg, navigationService)

	fmt.Println("bot is running")
	return handler.Start(ctx)
}
