package telegram

import (
    "context"
    "log"
    "strings"
    "time"

    tgadapter "cs-smokes-bot/internal/adapters/telegram"
    "cs-smokes-bot/internal/usecase"
)

type Handler struct {
    client *tgadapter.Client
    nav    *usecase.NavigationService
}

func NewHandler(client *tgadapter.Client, nav *usecase.NavigationService) *Handler {
    return &Handler{client: client, nav: nav}
}

func (h *Handler) Start(ctx context.Context) error {
    var offset int64
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        updates, err := h.client.GetUpdates(offset, 25)
        if err != nil {
            log.Printf("get updates error: %v", err)
            time.Sleep(2 * time.Second)
            continue
        }

        for _, update := range updates {
            offset = update.UpdateID + 1
            h.handleUpdate(update)
        }
    }
}

func (h *Handler) handleUpdate(update tgadapter.Update) {
    if update.Message != nil && update.Message.Text == "/start" {
        if err := h.nav.Start(update.Message.From.ID, update.Message.Chat.ID); err != nil {
            log.Printf("start handler error: %v", err)
        }
        return
    }

    if update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, "nav:") {
        if err := h.nav.HandleCallback(update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.ID, update.CallbackQuery.Data); err != nil {
            log.Printf("callback handler error: %v", err)
        }
    }
}
