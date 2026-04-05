package domain

type Screen string
type Mode string

const (
    ScreenMain      Screen = "main"
    ScreenHelp      Screen = "help"
    ScreenMaps      Screen = "maps"
    ScreenMapModes  Screen = "map_modes"
    ScreenZones     Screen = "zones"
    ScreenTargets   Screen = "targets"
    ScreenSmokeList Screen = "smoke_list"
    ScreenSmokeCard Screen = "smoke_card"
)

const (
    ModeZone   Mode = "zone"
    ModeTarget Mode = "target"
    ModeAll    Mode = "all"
)

type UserSession struct {
    UserID            int64
    ChatID            int64
    MenuMessageID     int
    ContentMessageID  *int
    CurrentScreen     Screen
    CurrentMapSlug    string
    CurrentMode       Mode
    CurrentZoneSlug   string
    CurrentTargetSlug string
    CurrentSmokeSlug  string
}
