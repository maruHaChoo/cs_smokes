package postgres

import (
	"context"
	"fmt"

	"cs-smokes-bot/internal/domain"
	"cs-smokes-bot/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SmokeRepository struct {
	pool *pgxpool.Pool
}

func NewSmokeRepository(pool *pgxpool.Pool) ports.SmokeRepository {
	return &SmokeRepository{
		pool: pool,
	}
}

func (r *SmokeRepository) GetMaps() ([]domain.GameMap, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT slug, title
		FROM maps
		ORDER BY title
	`)
	if err != nil {
		return nil, fmt.Errorf("get maps: %w", err)
	}
	defer rows.Close()

	var maps []domain.GameMap
	for rows.Next() {
		var item domain.GameMap
		if err := rows.Scan(&item.Slug, &item.Title); err != nil {
			return nil, fmt.Errorf("scan map: %w", err)
		}
		maps = append(maps, item)
	}

	return maps, rows.Err()
}

func (r *SmokeRepository) GetZonesByMap(mapSlug string) ([]string, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT DISTINCT zone_slug
		FROM smokes
		WHERE map_slug = $1 AND is_active = TRUE
		ORDER BY zone_slug
	`, mapSlug)
	if err != nil {
		return nil, fmt.Errorf("get zones by map: %w", err)
	}
	defer rows.Close()

	var zones []string
	for rows.Next() {
		var zone string
		if err := rows.Scan(&zone); err != nil {
			return nil, fmt.Errorf("scan zone: %w", err)
		}
		zones = append(zones, zone)
	}

	return zones, rows.Err()
}

func (r *SmokeRepository) GetTargetsByMap(mapSlug string) ([]string, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT DISTINCT target_slug
		FROM smokes
		WHERE map_slug = $1 AND is_active = TRUE
		ORDER BY target_slug
	`, mapSlug)
	if err != nil {
		return nil, fmt.Errorf("get targets by map: %w", err)
	}
	defer rows.Close()

	var targets []string
	for rows.Next() {
		var target string
		if err := rows.Scan(&target); err != nil {
			return nil, fmt.Errorf("scan target: %w", err)
		}
		targets = append(targets, target)
	}

	return targets, rows.Err()
}

func (r *SmokeRepository) ListByMap(mapSlug string) ([]domain.Smoke, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT
			id,
			map_slug,
			zone_slug,
			target_slug,
			slug,
			title,
			from_position,
			to_position,
			description,
			video_file_id,
			is_active
		FROM smokes
		WHERE map_slug = $1 AND is_active = TRUE
		ORDER BY title
	`, mapSlug)
	if err != nil {
		return nil, fmt.Errorf("list smokes by map: %w", err)
	}
	defer rows.Close()

	return scanSmokeRows(rows)
}

func (r *SmokeRepository) ListByMapAndZone(mapSlug, zoneSlug string) ([]domain.Smoke, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT
			id,
			map_slug,
			zone_slug,
			target_slug,
			slug,
			title,
			from_position,
			to_position,
			description,
			video_file_id,
			is_active
		FROM smokes
		WHERE map_slug = $1 AND zone_slug = $2 AND is_active = TRUE
		ORDER BY title
	`, mapSlug, zoneSlug)
	if err != nil {
		return nil, fmt.Errorf("list smokes by zone: %w", err)
	}
	defer rows.Close()

	return scanSmokeRows(rows)
}

func (r *SmokeRepository) ListByMapAndTarget(mapSlug, targetSlug string) ([]domain.Smoke, error) {
	rows, err := r.pool.Query(context.Background(), `
		SELECT
			id,
			map_slug,
			zone_slug,
			target_slug,
			slug,
			title,
			from_position,
			to_position,
			description,
			video_file_id,
			is_active
		FROM smokes
		WHERE map_slug = $1 AND target_slug = $2 AND is_active = TRUE
		ORDER BY title
	`, mapSlug, targetSlug)
	if err != nil {
		return nil, fmt.Errorf("list smokes by target: %w", err)
	}
	defer rows.Close()

	return scanSmokeRows(rows)
}

func (r *SmokeRepository) GetBySlug(slug string) (*domain.Smoke, error) {
	row := r.pool.QueryRow(context.Background(), `
		SELECT
			id,
			map_slug,
			zone_slug,
			target_slug,
			slug,
			title,
			from_position,
			to_position,
			description,
			video_file_id,
			is_active
		FROM smokes
		WHERE slug = $1 AND is_active = TRUE
	`, slug)

	var smoke domain.Smoke
	err := row.Scan(
		&smoke.ID,
		&smoke.MapSlug,
		&smoke.ZoneSlug,
		&smoke.TargetSlug,
		&smoke.Slug,
		&smoke.Title,
		&smoke.FromPosition,
		&smoke.ToPosition,
		&smoke.Description,
		&smoke.VideoFileID,
		&smoke.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("get smoke by slug: %w", err)
	}

	return &smoke, nil
}

type smokeRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func scanSmokeRows(rows smokeRows) ([]domain.Smoke, error) {
	var smokes []domain.Smoke

	for rows.Next() {
		var smoke domain.Smoke
		if err := rows.Scan(
			&smoke.ID,
			&smoke.MapSlug,
			&smoke.ZoneSlug,
			&smoke.TargetSlug,
			&smoke.Slug,
			&smoke.Title,
			&smoke.FromPosition,
			&smoke.ToPosition,
			&smoke.Description,
			&smoke.VideoFileID,
			&smoke.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan smoke: %w", err)
		}
		smokes = append(smokes, smoke)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return smokes, nil
}
