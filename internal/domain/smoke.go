package domain

type Smoke struct {
    ID           int64
    MapSlug      string
    ZoneSlug     string
    TargetSlug   string
    Slug         string
    Title        string
    FromPosition string
    ToPosition   string
    Description  string
    VideoFileID  string
    IsActive     bool
}
