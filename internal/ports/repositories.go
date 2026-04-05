package ports

import "cs-smokes-bot/internal/domain"

type SessionRepository interface {
    GetOrCreate(userID, chatID int64) (*domain.UserSession, error)
    Save(session *domain.UserSession) error
}

type SmokeRepository interface {
    GetMaps() ([]domain.GameMap, error)
    GetZonesByMap(mapSlug string) ([]string, error)
    GetTargetsByMap(mapSlug string) ([]string, error)
    ListByMap(mapSlug string) ([]domain.Smoke, error)
    ListByMapAndZone(mapSlug, zoneSlug string) ([]domain.Smoke, error)
    ListByMapAndTarget(mapSlug, targetSlug string) ([]domain.Smoke, error)
    GetBySlug(slug string) (*domain.Smoke, error)
}
