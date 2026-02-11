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

func (b *Bot) StartWebhook() {

}

func (b *Bot) Listener(){
	
}
var cfg = configs.Config{Port: 1023}

func Setup() {
	mux := http.NewServeMux()
	mux.HandleFunc("", WebHookHandler)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
func WebHookHandler(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, 30)
	r.Body.Read(data)
	print(string(data))
}
