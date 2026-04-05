package main

import (
    "context"
    "log"

    "cs-smokes-bot/internal/app"
)

func main() {
    ctx := context.Background()

    if err := app.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
