package memory

import (
    "fmt"
    "sort"

    "cs-smokes-bot/internal/domain"
)

type SmokeRepository struct {
    maps   []domain.GameMap
    smokes []domain.Smoke
}

func NewSmokeRepository() *SmokeRepository {
    return &SmokeRepository{
        maps: []domain.GameMap{
            {Slug: "mirage", Title: "Mirage"},
            {Slug: "inferno", Title: "Inferno"},
        },
        smokes: []domain.Smoke{
            {ID: 1, MapSlug: "mirage", ZoneSlug: "mid", TargetSlug: "window", Slug: "window_from_t_spawn", Title: "Window from T-spawn", FromPosition: "T-spawn", ToPosition: "Window", Description: "Быстрый стандартный смок в окно.", VideoFileID: "REPLACE_WITH_REAL_FILE_ID", IsActive: true},
            {ID: 2, MapSlug: "mirage", ZoneSlug: "a", TargetSlug: "jungle", Slug: "jungle_from_ramp", Title: "Jungle from Ramp", FromPosition: "Ramp", ToPosition: "Jungle", Description: "Классический выходной смок в джангл.", VideoFileID: "REPLACE_WITH_REAL_FILE_ID", IsActive: true},
            {ID: 3, MapSlug: "mirage", ZoneSlug: "a", TargetSlug: "ct", Slug: "ct_from_tetris", Title: "CT from Tetris", FromPosition: "Tetris", ToPosition: "CT", Description: "Выходной смок в КТ на A site.", VideoFileID: "REPLACE_WITH_REAL_FILE_ID", IsActive: true},
            {ID: 4, MapSlug: "inferno", ZoneSlug: "b", TargetSlug: "coffins", Slug: "coffins_from_banana", Title: "Coffins from Banana", FromPosition: "Banana", ToPosition: "Coffins", Description: "Смок в гробы с банана.", VideoFileID: "REPLACE_WITH_REAL_FILE_ID", IsActive: true},
            {ID: 5, MapSlug: "inferno", ZoneSlug: "a", TargetSlug: "library", Slug: "library_from_second_mid", Title: "Library from Second Mid", FromPosition: "Second Mid", ToPosition: "Library", Description: "Смок в библиотеку для выхода на A.", VideoFileID: "REPLACE_WITH_REAL_FILE_ID", IsActive: true},
        },
    }
}

func (r *SmokeRepository) GetMaps() ([]domain.GameMap, error) {
    out := make([]domain.GameMap, len(r.maps))
    copy(out, r.maps)
    return out, nil
}

func (r *SmokeRepository) GetZonesByMap(mapSlug string) ([]string, error) {
    uniq := map[string]struct{}{}
    for _, smoke := range r.smokes {
        if smoke.MapSlug == mapSlug && smoke.IsActive { uniq[smoke.ZoneSlug] = struct{}{} }
    }
    out := keys(uniq)
    sort.Strings(out)
    return out, nil
}

func (r *SmokeRepository) GetTargetsByMap(mapSlug string) ([]string, error) {
    uniq := map[string]struct{}{}
    for _, smoke := range r.smokes {
        if smoke.MapSlug == mapSlug && smoke.IsActive { uniq[smoke.TargetSlug] = struct{}{} }
    }
    out := keys(uniq)
    sort.Strings(out)
    return out, nil
}

func (r *SmokeRepository) ListByMap(mapSlug string) ([]domain.Smoke, error) {
    out := make([]domain.Smoke, 0)
    for _, smoke := range r.smokes {
        if smoke.MapSlug == mapSlug && smoke.IsActive { out = append(out, smoke) }
    }
    sortSmokes(out)
    return out, nil
}

func (r *SmokeRepository) ListByMapAndZone(mapSlug, zoneSlug string) ([]domain.Smoke, error) {
    out := make([]domain.Smoke, 0)
    for _, smoke := range r.smokes {
        if smoke.MapSlug == mapSlug && smoke.ZoneSlug == zoneSlug && smoke.IsActive { out = append(out, smoke) }
    }
    sortSmokes(out)
    return out, nil
}

func (r *SmokeRepository) ListByMapAndTarget(mapSlug, targetSlug string) ([]domain.Smoke, error) {
    out := make([]domain.Smoke, 0)
    for _, smoke := range r.smokes {
        if smoke.MapSlug == mapSlug && smoke.TargetSlug == targetSlug && smoke.IsActive { out = append(out, smoke) }
    }
    sortSmokes(out)
    return out, nil
}

func (r *SmokeRepository) GetBySlug(slug string) (*domain.Smoke, error) {
    for _, smoke := range r.smokes {
        if smoke.Slug == slug && smoke.IsActive {
            cp := smoke
            return &cp, nil
        }
    }
    return nil, fmt.Errorf("smoke not found: %s", slug)
}

func keys(m map[string]struct{}) []string {
    out := make([]string, 0, len(m))
    for k := range m { out = append(out, k) }
    return out
}

func sortSmokes(smokes []domain.Smoke) {
    sort.Slice(smokes, func(i, j int) bool { return smokes[i].Title < smokes[j].Title })
}
