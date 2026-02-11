package telego

type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Data    string   `json:"data"`
	Message *Message `json:"message"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
	URL          string `json:"url,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}
