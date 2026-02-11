package telego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Bot struct {
	Token            string
	BaseURL          string
	Offset           int
	Client           *http.Client
	Handlers         map[string]func(*Context)
	MessageHandlers  []func(*Context) bool
	CallbackHandlers map[string]func(*Context) bool
	Logger           *log.Logger
	Debug            bool
}

func NewBot(token string) *Bot {
	return &Bot{
		Token:            token,
		BaseURL:          fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Client:           &http.Client{Timeout: 10 * time.Second},
		Handlers:         make(map[string]func(*Context)),
		MessageHandlers:  make([]func(*Context) bool, 0),
		CallbackHandlers: make(map[string]func(*Context) bool),
		Logger:           log.New(log.Writer(), "[telego] ", log.LstdFlags),
		Debug:            false,
	}
}

func (b *Bot) HandleCommand(cmd string, handler func(*Context)) {
	b.Handlers[cmd] = handler
}

func (b *Bot) HandleMessage(handler func(*Context) bool) {
	b.MessageHandlers = append(b.MessageHandlers, handler)
}

func (b *Bot) HandleCallback(data string, handler func(*Context) bool) {
	b.CallbackHandlers[data] = handler
}

func (b *Bot) Start() {
	b.Logger.Println("Bot başlatılıyor...")
	for {
		updates, err := b.getUpdates()
		if err != nil {
			b.Logger.Println("Polling hatası:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, update := range updates {
			if update.Message != nil {
				b.handleMessage(update.Message)
			}
			if update.CallbackQuery != nil {
				b.handleCallback(update.CallbackQuery)
			}
		}
	}
}

func (b *Bot) handleMessage(msg *Message) {
	ctx := &Context{
		Bot:     b,
		ChatID:  msg.Chat.ID,
		Message: msg,
	}

	if b.Debug {
		b.Logger.Printf("Mesaj alındı [%d]: %s", msg.Chat.ID, msg.Text)
	}

	// Komut kontrolü
	if strings.HasPrefix(msg.Text, "/") {
		parts := strings.Fields(msg.Text)
		cmd := parts[0]
		ctx.Args = parts[1:]

		if handler, ok := b.Handlers[cmd]; ok {
			handler(ctx)
			return
		}
	}

	// Genel mesaj handler'ları
	for _, handler := range b.MessageHandlers {
		if handler(ctx) {
			return // Handler True dönürürse diğerleri çalışmasın
		}
	}
}

func (b *Bot) handleCallback(cb *CallbackQuery) {
	if handler, ok := b.CallbackHandlers[cb.Data]; ok {
		ctx := &Context{
			Bot:           b,
			CallbackQuery: cb,
			ChatID:        cb.Message.Chat.ID,
			Message:       cb.Message,
		}
		if b.Debug {
			b.Logger.Printf("Callback alındı [%d]: %s", cb.Message.Chat.ID, cb.Data)
		}
		handler(ctx)
	}
}

func (b *Bot) getUpdates() ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d", b.BaseURL, b.Offset)
	resp, err := b.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Result) > 0 {
		b.Offset = result.Result[len(result.Result)-1].UpdateID + 1
	}

	return result.Result, nil
}

// ============ Public API Methods ============

// SendMessage bot tarafından bir sohbete metin mesajı gönderir
func (b *Bot) SendMessage(chatID int64, text string) error {
	endpoint := fmt.Sprintf("%s/sendMessage", b.BaseURL)
	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("%d", chatID))
	data.Set("text", text)

	resp, err := b.Client.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API hatası: %d", resp.StatusCode)
	}

	return nil
}

// SendMessageWithKeyboard inline keyboard ile mesaj gönderir
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard *InlineKeyboardMarkup) error {
	endpoint := fmt.Sprintf("%s/sendMessage", b.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": keyboard,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := b.Client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API hatası: %d", resp.StatusCode)
	}

	return nil
}

// EditMessage bir mesajı düzenler
func (b *Bot) EditMessage(chatID int64, messageID int, text string) error {
	endpoint := fmt.Sprintf("%s/editMessageText", b.BaseURL)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := b.Client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API hatası: %d", resp.StatusCode)
	}

	return nil
}

// DeleteMessage bir mesajı siler
func (b *Bot) DeleteMessage(chatID int64, messageID int) error {
	endpoint := fmt.Sprintf("%s/deleteMessage", b.BaseURL)
	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("%d", chatID))
	data.Set("message_id", fmt.Sprintf("%d", messageID))

	resp, err := b.Client.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API hatası: %d", resp.StatusCode)
	}

	return nil
}
