# CS Smokes Bot

MVP Telegram bot on Go with single-menu-message navigation and auto-cleanup of smoke videos.

## What is implemented
- one permanent menu message per user
- smoke navigation:
  - by map
  - by mode: zones / targets / all smokes
- smoke card screen
- previous smoke video is deleted when:
  - another smoke is opened
  - user leaves smoke card via back/menu
  - user switches branch
- scalable architecture with separated domain / usecase / ports / adapters
- in-memory repositories for now

## Run
1. Export bot token:

```bash
export BOT_TOKEN=your_token_here
```

2. Start:

```bash
go run ./cmd/bot
```

Replace `REPLACE_WITH_REAL_FILE_ID` in the memory smoke repository with real Telegram `file_id` values.


## Docker
1. Create env file:

```bash
cp .env.example .env
```

2. Fill `BOT_TOKEN` in `.env`.

3. Run:

```bash
docker compose up -d --build
```

## Migrations
Prepared files:
- `migrations/000001_init.up.sql`
- `migrations/000001_init.down.sql`

They create base tables for:
- users
- user_sessions
- maps
- smokes

Right now the bot still uses in-memory repositories. Migrations, Dockerfile and Compose are ready for the next step: switching storage to Postgres.
