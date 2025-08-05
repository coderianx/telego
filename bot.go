package telego

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Bot struct {
	Token    string
	BaseURL  string
	Offset   int
	Client   *http.Client
	Handlers map[string]func(*Context)
}

func NewBot(token string) *Bot {
	return &Bot{
		Token:    token,
		BaseURL:  fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Client:   &http.Client{Timeout: 10 * time.Second},
		Handlers: make(map[string]func(*Context)),
	}
}

func (b *Bot) HandleCommand(cmd string, handler func(*Context)) {
	b.Handlers[cmd] = handler
}

func (b *Bot) Start() {
	for {
		updates, err := b.getUpdates()
		if err != nil {
			fmt.Println("Polling hatası:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, update := range updates {
			msg := update.Message
			handler, ok := b.Handlers[msg.Text]
			if ok {
				ctx := &Context{
					Bot:     b,
					ChatID:  msg.Chat.ID,
					Message: msg,
				}
				handler(ctx)
			}
		}
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
