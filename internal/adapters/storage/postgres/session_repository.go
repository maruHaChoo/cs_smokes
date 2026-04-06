package postgres

import (
	"context"
	"errors"
	"fmt"

	"cs-smokes-bot/internal/domain"
	"cs-smokes-bot/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) ports.SessionRepository {
	return &SessionRepository{
		pool: pool,
	}
}

func (r *SessionRepository) GetOrCreate(userID, chatID int64) (*domain.UserSession, error) {
	ctx := context.Background()

	session, err := r.getByUserID(ctx, userID)
	if err == nil {
		if session.ChatID != chatID {
			session.ChatID = chatID
			if err := r.Save(session); err != nil {
				return nil, err
			}
		}
		return session, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("get session: %w", err)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO user_sessions (
			user_id,
			chat_id,
			menu_message_id,
			content_message_id,
			current_screen,
			current_map_slug,
			current_mode,
			current_zone_slug,
			current_target_slug,
			current_smoke_slug
		)
		VALUES ($1, $2, 0, NULL, '', '', '', '', '', '')
	`, userID, chatID)
	if err != nil {
		return nil, fmt.Errorf("insert session: %w", err)
	}

	return &domain.UserSession{
		UserID:            userID,
		ChatID:            chatID,
		MenuMessageID:     0,
		ContentMessageID:  nil,
		CurrentScreen:     "",
		CurrentMapSlug:    "",
		CurrentMode:       "",
		CurrentZoneSlug:   "",
		CurrentTargetSlug: "",
		CurrentSmokeSlug:  "",
	}, nil
}

func (r *SessionRepository) Save(session *domain.UserSession) error {
	ctx := context.Background()

	_, err := r.pool.Exec(ctx, `
		UPDATE user_sessions
		SET
			chat_id = $2,
			menu_message_id = $3,
			content_message_id = $4,
			current_screen = $5,
			current_map_slug = $6,
			current_mode = $7,
			current_zone_slug = $8,
			current_target_slug = $9,
			current_smoke_slug = $10
		WHERE user_id = $1
	`,
		session.UserID,
		session.ChatID,
		session.MenuMessageID,
		session.ContentMessageID,
		string(session.CurrentScreen),
		session.CurrentMapSlug,
		string(session.CurrentMode),
		session.CurrentZoneSlug,
		session.CurrentTargetSlug,
		session.CurrentSmokeSlug,
	)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}

	return nil
}

func (r *SessionRepository) getByUserID(ctx context.Context, userID int64) (*domain.UserSession, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			user_id,
			chat_id,
			menu_message_id,
			content_message_id,
			current_screen,
			current_map_slug,
			current_mode,
			current_zone_slug,
			current_target_slug,
			current_smoke_slug
		FROM user_sessions
		WHERE user_id = $1
	`, userID)

	var session domain.UserSession
	var currentScreen string
	var currentMode string

	err := row.Scan(
		&session.UserID,
		&session.ChatID,
		&session.MenuMessageID,
		&session.ContentMessageID,
		&currentScreen,
		&session.CurrentMapSlug,
		&currentMode,
		&session.CurrentZoneSlug,
		&session.CurrentTargetSlug,
		&session.CurrentSmokeSlug,
	)
	if err != nil {
		return nil, err
	}

	session.CurrentScreen = domain.Screen(currentScreen)
	session.CurrentMode = domain.Mode(currentMode)

	return &session, nil
}
