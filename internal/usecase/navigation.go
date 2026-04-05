package usecase

import (
    "fmt"
    "strings"

    "cs-smokes-bot/internal/domain"
    "cs-smokes-bot/internal/ports"
)

type NavigationService struct {
    sessions ports.SessionRepository
    smokes   ports.SmokeRepository
    tg       ports.TelegramGateway
}

func NewNavigationService(sessions ports.SessionRepository, smokes ports.SmokeRepository, tg ports.TelegramGateway) *NavigationService {
    return &NavigationService{sessions: sessions, smokes: smokes, tg: tg}
}

func (s *NavigationService) Start(userID, chatID int64) error {
    session, err := s.sessions.GetOrCreate(userID, chatID)
    if err != nil {
        return err
    }

    text, rows := s.renderMain()
    if session.MenuMessageID == 0 {
        menuID, err := s.tg.SendMenu(chatID, text, rows)
        if err != nil {
            return err
        }
        session.MenuMessageID = menuID
    } else {
        if err := s.cleanupContent(session); err != nil {
            return err
        }
        if err := s.tg.EditMenu(chatID, session.MenuMessageID, text, rows); err != nil {
            return err
        }
    }

    session.CurrentScreen = domain.ScreenMain
    session.CurrentMapSlug = ""
    session.CurrentMode = ""
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) HandleCallback(userID, chatID int64, callbackID, data string) error {
    session, err := s.sessions.GetOrCreate(userID, chatID)
    if err != nil {
        return err
    }
    defer func() { _ = s.tg.AnswerCallback(callbackID) }()

    parts := strings.Split(data, ":")
    if len(parts) < 2 || parts[0] != "nav" {
        return nil
    }

    switch parts[1] {
    case "main":
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showMain(session)
    case "help":
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showHelp(session)
    case "maps":
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showMaps(session)
    case "map_modes":
        if len(parts) < 3 { return nil }
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showMapModes(session, parts[2])
    case "zones":
        if len(parts) < 3 { return nil }
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showZones(session, parts[2])
    case "targets":
        if len(parts) < 3 { return nil }
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showTargets(session, parts[2])
    case "all_smokes":
        if len(parts) < 3 { return nil }
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showAllSmokes(session, parts[2])
    case "smoke_list_zone":
        if len(parts) < 4 { return nil }
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showSmokeListByZone(session, parts[2], parts[3])
    case "smoke_list_target":
        if len(parts) < 4 { return nil }
        if err := s.cleanupContent(session); err != nil { return err }
        return s.showSmokeListByTarget(session, parts[2], parts[3])
    case "smoke_card":
        if len(parts) < 3 { return nil }
        return s.showSmokeCard(session, parts[2])
    case "back_from_card":
        if err := s.cleanupContent(session); err != nil { return err }
        return s.backFromCard(session)
    default:
        return nil
    }
}

func (s *NavigationService) showMain(session *domain.UserSession) error {
    text, rows := s.renderMain()
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenMain
    session.CurrentMapSlug = ""
    session.CurrentMode = ""
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showHelp(session *domain.UserSession) error {
    text := "ℹ️ *Как пользоваться*\n\n1. Выбери карту\n2. Выбери режим: по зонам / по точкам / все смоки\n3. Открой нужный смок\n4. Бот отправит видео и очистит его при выходе назад"
    rows := [][]ports.Button{{{Text: "🎯 Смоки", Data: "nav:maps"}}, {{Text: "⌂ В меню", Data: "nav:main"}}}
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenHelp
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showMaps(session *domain.UserSession) error {
    maps, err := s.smokes.GetMaps()
    if err != nil { return err }
    text := "🗺 *Выбери карту*"
    rows := make([][]ports.Button, 0, len(maps)+1)
    for _, gameMap := range maps {
        rows = append(rows, []ports.Button{{Text: gameMap.Title, Data: fmt.Sprintf("nav:map_modes:%s", gameMap.Slug)}})
    }
    rows = append(rows, []ports.Button{{Text: "⌂ В меню", Data: "nav:main"}})
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenMaps
    session.CurrentMapSlug = ""
    session.CurrentMode = ""
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showMapModes(session *domain.UserSession, mapSlug string) error {
    text := fmt.Sprintf("🗺 *%s*\n\nВыбери режим поиска смоков.", titleCase(mapSlug))
    rows := [][]ports.Button{
        {{Text: "🎯 По зонам", Data: fmt.Sprintf("nav:zones:%s", mapSlug)}},
        {{Text: "📍 По точкам", Data: fmt.Sprintf("nav:targets:%s", mapSlug)}},
        {{Text: "📋 Все смоки", Data: fmt.Sprintf("nav:all_smokes:%s", mapSlug)}},
        {{Text: "← Назад", Data: "nav:maps"}, {Text: "⌂ В меню", Data: "nav:main"}},
    }
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenMapModes
    session.CurrentMapSlug = mapSlug
    session.CurrentMode = ""
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showZones(session *domain.UserSession, mapSlug string) error {
    zones, err := s.smokes.GetZonesByMap(mapSlug)
    if err != nil { return err }
    text := fmt.Sprintf("🎯 *%s / По зонам*\n\nВыбери зону.", titleCase(mapSlug))
    rows := make([][]ports.Button, 0, len(zones)+1)
    for _, zone := range zones {
        rows = append(rows, []ports.Button{{Text: strings.ToUpper(zone), Data: fmt.Sprintf("nav:smoke_list_zone:%s:%s", mapSlug, zone)}})
    }
    rows = append(rows, []ports.Button{{Text: "← Назад", Data: fmt.Sprintf("nav:map_modes:%s", mapSlug)}, {Text: "⌂ В меню", Data: "nav:main"}})
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenZones
    session.CurrentMapSlug = mapSlug
    session.CurrentMode = domain.ModeZone
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showTargets(session *domain.UserSession, mapSlug string) error {
    targets, err := s.smokes.GetTargetsByMap(mapSlug)
    if err != nil { return err }
    text := fmt.Sprintf("📍 *%s / По точкам*\n\nВыбери целевую точку.", titleCase(mapSlug))
    rows := make([][]ports.Button, 0, len(targets)+1)
    for _, target := range targets {
        rows = append(rows, []ports.Button{{Text: titleCase(target), Data: fmt.Sprintf("nav:smoke_list_target:%s:%s", mapSlug, target)}})
    }
    rows = append(rows, []ports.Button{{Text: "← Назад", Data: fmt.Sprintf("nav:map_modes:%s", mapSlug)}, {Text: "⌂ В меню", Data: "nav:main"}})
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenTargets
    session.CurrentMapSlug = mapSlug
    session.CurrentMode = domain.ModeTarget
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showAllSmokes(session *domain.UserSession, mapSlug string) error {
    smokes, err := s.smokes.ListByMap(mapSlug)
    if err != nil { return err }
    text := fmt.Sprintf("📋 *%s / Все смоки*\n\nВыбери смок.", titleCase(mapSlug))
    rows := s.buildSmokeRows(smokes)
    rows = append(rows, []ports.Button{{Text: "← Назад", Data: fmt.Sprintf("nav:map_modes:%s", mapSlug)}, {Text: "⌂ В меню", Data: "nav:main"}})
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenSmokeList
    session.CurrentMapSlug = mapSlug
    session.CurrentMode = domain.ModeAll
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showSmokeListByZone(session *domain.UserSession, mapSlug, zoneSlug string) error {
    smokes, err := s.smokes.ListByMapAndZone(mapSlug, zoneSlug)
    if err != nil { return err }
    text := fmt.Sprintf("🎯 *%s / %s*\n\nВыбери смок.", titleCase(mapSlug), strings.ToUpper(zoneSlug))
    rows := s.buildSmokeRows(smokes)
    rows = append(rows, []ports.Button{{Text: "← Назад", Data: fmt.Sprintf("nav:zones:%s", mapSlug)}, {Text: "⌂ В меню", Data: "nav:main"}})
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenSmokeList
    session.CurrentMapSlug = mapSlug
    session.CurrentMode = domain.ModeZone
    session.CurrentZoneSlug = zoneSlug
    session.CurrentTargetSlug = ""
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showSmokeListByTarget(session *domain.UserSession, mapSlug, targetSlug string) error {
    smokes, err := s.smokes.ListByMapAndTarget(mapSlug, targetSlug)
    if err != nil { return err }
    text := fmt.Sprintf("📍 *%s / %s*\n\nВыбери смок.", titleCase(mapSlug), titleCase(targetSlug))
    rows := s.buildSmokeRows(smokes)
    rows = append(rows, []ports.Button{{Text: "← Назад", Data: fmt.Sprintf("nav:targets:%s", mapSlug)}, {Text: "⌂ В меню", Data: "nav:main"}})
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }
    session.CurrentScreen = domain.ScreenSmokeList
    session.CurrentMapSlug = mapSlug
    session.CurrentMode = domain.ModeTarget
    session.CurrentZoneSlug = ""
    session.CurrentTargetSlug = targetSlug
    session.CurrentSmokeSlug = ""
    return s.sessions.Save(session)
}

func (s *NavigationService) showSmokeCard(session *domain.UserSession, smokeSlug string) error {
    smoke, err := s.smokes.GetBySlug(smokeSlug)
    if err != nil { return err }
    if err := s.cleanupContent(session); err != nil { return err }

    text := fmt.Sprintf("🎬 *%s*\n\n🗺 Карта: %s\n📍 Откуда: %s\n🎯 Куда: %s\n📝 %s",
        smoke.Title, titleCase(smoke.MapSlug), smoke.FromPosition, smoke.ToPosition, smoke.Description,
    )
    rows := [][]ports.Button{{{Text: "← Назад", Data: "nav:back_from_card"}, {Text: "⌂ В меню", Data: "nav:main"}}}
    if err := s.tg.EditMenu(session.ChatID, session.MenuMessageID, text, rows); err != nil { return err }

    videoMessageID, err := s.tg.SendVideo(session.ChatID, smoke.VideoFileID, smoke.Title)
    if err != nil { return err }

    session.ContentMessageID = &videoMessageID
    session.CurrentScreen = domain.ScreenSmokeCard
    session.CurrentSmokeSlug = smokeSlug
    session.CurrentMapSlug = smoke.MapSlug
    session.CurrentZoneSlug = smoke.ZoneSlug
    session.CurrentTargetSlug = smoke.TargetSlug
    return s.sessions.Save(session)
}

func (s *NavigationService) backFromCard(session *domain.UserSession) error {
    switch session.CurrentMode {
    case domain.ModeZone:
        return s.showSmokeListByZone(session, session.CurrentMapSlug, session.CurrentZoneSlug)
    case domain.ModeTarget:
        return s.showSmokeListByTarget(session, session.CurrentMapSlug, session.CurrentTargetSlug)
    case domain.ModeAll:
        return s.showAllSmokes(session, session.CurrentMapSlug)
    default:
        return s.showMaps(session)
    }
}

func (s *NavigationService) cleanupContent(session *domain.UserSession) error {
    if session.ContentMessageID == nil { return nil }
    _ = s.tg.DeleteMessage(session.ChatID, *session.ContentMessageID)
    session.ContentMessageID = nil
    return s.sessions.Save(session)
}

func (s *NavigationService) renderMain() (string, [][]ports.Button) {
    text := "🔥 *CS Smokes Bot*\n\nВыбирай карту, ищи нужный смок и смотри видео без мусора в чате.\n\nПоддерживается единый экран меню и автоочистка прошлого видео."
    rows := [][]ports.Button{
        {{Text: "🎯 Смоки", Data: "nav:maps"}},
        {{Text: "ℹ️ Помощь", Data: "nav:help"}},
    }
    return text, rows
}

func (s *NavigationService) buildSmokeRows(smokes []domain.Smoke) [][]ports.Button {
    rows := make([][]ports.Button, 0, len(smokes))
    for _, smoke := range smokes {
        rows = append(rows, []ports.Button{{Text: smoke.Title, Data: fmt.Sprintf("nav:smoke_card:%s", smoke.Slug)}})
    }
    return rows
}

func titleCase(v string) string {
    if v == "" { return v }
    parts := strings.Split(v, "_")
    for i, part := range parts {
        if part == "" { continue }
        parts[i] = strings.ToUpper(part[:1]) + part[1:]
    }
    return strings.Join(parts, " ")
}
