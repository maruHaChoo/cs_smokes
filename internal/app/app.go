package app

import (
    "context"
    "fmt"

    memorystorage "cs-smokes-bot/internal/adapters/storage/memory"
    tgadapter "cs-smokes-bot/internal/adapters/telegram"
    "cs-smokes-bot/internal/config"
    tgtransport "cs-smokes-bot/internal/transport/telegram"
    "cs-smokes-bot/internal/usecase"
)

func Run(ctx context.Context) error {
    cfg, err := config.Load()
    if err != nil { return err }

    tg, err := tgadapter.NewClient(cfg.BotToken)
    if err != nil { return err }

    sessionRepo := memorystorage.NewSessionRepository()
    smokeRepo := memorystorage.NewSmokeRepository()
    navigationService := usecase.NewNavigationService(sessionRepo, smokeRepo, tg)
    handler := tgtransport.NewHandler(tg, navigationService)

    fmt.Println("bot is running")
    return handler.Start(ctx)
}
