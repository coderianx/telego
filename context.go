package telego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Context struct {
	Bot           *Bot
	ChatID        int64
	Message       *Message
	CallbackQuery *CallbackQuery
	Args          []string
}

// ============ Helper Functions ============

// Text mesaj içeriğini döndürür
func (c *Context) Text() string {
	if c.Message != nil {
		return c.Message.Text
	}
	return ""
}

// UserName kullanıcı adını döndürür
func (c *Context) UserName() string {
	if c.Message != nil {
		return c.Message.From.Username
	}
	if c.CallbackQuery != nil {
		return c.CallbackQuery.From.Username
	}
	return ""
}

// FirstName kullanıcının adını döndürür
func (c *Context) FirstName() string {
	if c.Message != nil {
		return c.Message.From.FirstName
	}
	if c.CallbackQuery != nil {
		return c.CallbackQuery.From.FirstName
	}
	return ""
}

// UserID kullanıcı ID'sini döndürür
func (c *Context) UserID() int {
	if c.Message != nil {
		return c.Message.From.ID
	}
	if c.CallbackQuery != nil {
		return c.CallbackQuery.From.ID
	}
	return 0
}

// MessageID mesaj ID'sini döndürür
func (c *Context) MessageID() int {
	if c.Message != nil {
		return c.Message.MessageID
	}
	if c.CallbackQuery != nil && c.CallbackQuery.Message != nil {
		return c.CallbackQuery.Message.MessageID
	}
	return 0
}

// ============ Telegram API Methods ============

// SendMessage basit metin mesajı gönderir
func (c *Context) SendMessage(text string) error {
	return c.Bot.SendMessage(c.ChatID, text)
}

// SendMessageWithKeyboard inline keyboard ile mesaj gönderir
func (c *Context) SendMessageWithKeyboard(text string, keyboard *InlineKeyboardMarkup) error {
	endpoint := fmt.Sprintf("%s/sendMessage", c.Bot.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      c.ChatID,
		"text":         text,
		"reply_markup": keyboard,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Bot.Client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API hatası: %d", resp.StatusCode)
	}

	return nil
}

// EditMessage mesajı düzenler
func (c *Context) EditMessage(messageID int, text string) error {
	return c.Bot.EditMessage(c.ChatID, messageID, text)
}

// EditMessageWithKeyboard mesajı keyboard ile birlikte düzenler
func (c *Context) EditMessageWithKeyboard(messageID int, text string, keyboard *InlineKeyboardMarkup) error {
	endpoint := fmt.Sprintf("%s/editMessageText", c.Bot.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      c.ChatID,
		"message_id":   messageID,
		"text":         text,
		"reply_markup": keyboard,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := c.Bot.Client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API hatası: %d", resp.StatusCode)
	}

	return nil
}

// DeleteMessage mesajı siler
func (c *Context) DeleteMessage(messageID int) error {
	return c.Bot.DeleteMessage(c.ChatID, messageID)
}

// AnswerCallback callback sorgusu yanıtlar
func (c *Context) AnswerCallback(text string, alert bool) error {
	if c.CallbackQuery == nil {
		return fmt.Errorf("callback query yok")
	}
	endpoint := fmt.Sprintf("%s/answerCallbackQuery", c.Bot.BaseURL)

	data := url.Values{}
	data.Set("callback_query_id", c.CallbackQuery.ID)
	data.Set("text", text)
	if alert {
		data.Set("show_alert", "true")
	}

	resp, err := c.Bot.Client.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
