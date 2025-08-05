package telego

import (
	"fmt"
	"net/url"
)

type Context struct {
	Bot     *Bot
	ChatID  int64
	Message Message
}

func (c *Context) SendMessage(text string) error {
	endpoint := fmt.Sprintf("%s/sendMessage", c.Bot.BaseURL)
	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("%d", c.ChatID))
	data.Set("text", text)

	_, err := c.Bot.Client.PostForm(endpoint, data)
	return err
}
