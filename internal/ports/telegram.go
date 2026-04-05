package ports

type Button struct {
    Text string
    Data string
}

type TelegramGateway interface {
    SendMenu(chatID int64, text string, rows [][]Button) (int, error)
    EditMenu(chatID int64, messageID int, text string, rows [][]Button) error
    SendVideo(chatID int64, fileID, caption string) (int, error)
    DeleteMessage(chatID int64, messageID int) error
    AnswerCallback(callbackID string) error
}
