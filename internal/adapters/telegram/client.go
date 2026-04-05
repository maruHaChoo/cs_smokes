package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cs-smokes-bot/internal/ports"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("empty telegram token")
	}
	return &Client{
		baseURL:    fmt.Sprintf("https://api.telegram.org/bot%s", token),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *Client) SendMenu(chatID int64, text string, rows [][]ports.Button) (int, error) {
	payload := map[string]any{"chat_id": chatID, "text": text, "parse_mode": "Markdown", "reply_markup": toReplyMarkup(rows)}
	var response sendMessageResponse
	if err := c.call("sendMessage", payload, &response); err != nil {
		return 0, err
	}
	return response.Result.MessageID, nil
}

func (c *Client) EditMenu(chatID int64, messageID int, text string, rows [][]ports.Button) error {
	payload := map[string]any{"chat_id": chatID, "message_id": messageID, "text": text, "parse_mode": "Markdown", "reply_markup": toReplyMarkup(rows)}
	return c.call("editMessageText", payload, nil)
}

func (c *Client) SendVideo(chatID int64, fileID, caption string) (int, error) {
	payload := map[string]any{"chat_id": chatID, "video": fileID, "caption": caption}
	var response sendMessageResponse
	if err := c.call("sendVideo", payload, &response); err != nil {
		return 0, err
	}
	return response.Result.MessageID, nil
}

func (c *Client) DeleteMessage(chatID int64, messageID int) error {
	payload := map[string]any{"chat_id": chatID, "message_id": messageID}
	return c.call("deleteMessage", payload, nil)
}

func (c *Client) AnswerCallback(callbackID string) error {
	payload := map[string]any{"callback_query_id": callbackID}
	return c.call("answerCallbackQuery", payload, nil)
}

func (c *Client) GetUpdates(offset int64, timeoutSeconds int) ([]Update, error) {
	payload := map[string]any{"offset": offset, "timeout": timeoutSeconds}
	var response getUpdatesResponse
	if err := c.call("getUpdates", payload, &response); err != nil {
		return nil, err
	}
	return response.Result, nil
}

func (c *Client) call(method string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/"+method, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram returned status %d", resp.StatusCode)
	}

	raw := struct {
		OK          bool            `json:"ok"`
		Description string          `json:"description"`
		Result      json.RawMessage `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return fmt.Errorf("decode telegram response: %w", err)
	}
	if !raw.OK {
		return fmt.Errorf("telegram api error: %s", raw.Description)
	}
	if out == nil {
		return nil
	}

	wrapperBody, err := json.Marshal(map[string]any{"result": raw.Result})
	if err != nil {
		return fmt.Errorf("marshal wrapper response: %w", err)
	}
	if err := json.Unmarshal(wrapperBody, out); err != nil {
		return fmt.Errorf("unmarshal result body: %w", err)
	}
	return nil
}

func toReplyMarkup(rows [][]ports.Button) map[string]any {
	keyboard := make([][]map[string]string, 0, len(rows))
	for _, row := range rows {
		items := make([]map[string]string, 0, len(row))
		for _, btn := range row {
			items = append(items, map[string]string{"text": btn.Text, "callback_data": btn.Data})
		}
		keyboard = append(keyboard, items)
	}
	return map[string]any{"inline_keyboard": keyboard}
}

type sendMessageResponse struct {
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

type getUpdatesResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text,omitempty"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
}

type CallbackQuery struct {
	ID      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type Chat struct {
	ID int64 `json:"id"`
}
type User struct {
	ID int64 `json:"id"`
}
