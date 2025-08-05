# mytelegrambot

**A lightweight Telegram Bot library written in Go.**

---

## ✨ Features

- ✅ Easy bot initialization with `NewBot("TOKEN")`
- ✅ Command handling via `HandleCommand("/command", func(ctx *Context))`
- ✅ Simple message sending with `ctx.SendMessage("...")`
- ✅ Auto polling with `Start()`

---

## 🚀 Installation

```bash
go get github.com/coderianx/telego
```

---

## 🔰 Example Usage

```go
package main

import (
    "log"
    "github.com/coderianx/telego"
)

func main() {
    bot := telego.NewBot("YOUR_BOT_TOKEN")

    bot.HandleCommand("/start", func(ctx *telego.Context) {
        ctx.SendMessage("Welcome!")
    })

    bot.HandleCommand("/ping", func(ctx *telego.Context) {
        ctx.SendMessage("Pong 🏓")
    })

    log.Println("Bot is running...")
    bot.Start()
}
```

---

## ✅ Public API

| Function                          | Description |
|-----------------------------------|-------------|
| `NewBot(token string)`            | Creates a new bot instance |
| `HandleCommand(cmd string, fn)`   | Assigns a handler to a command |
| `ctx.SendMessage(text string)`    | Sends a message to the chat |
| `bot.Start()`              | Starts polling for updates |

---

## 🛠️ Roadmap

- [ ] Helper functions like `ctx.Text()`
- [ ] Inline keyboard support
- [ ] Sending photos/documents
- [ ] Webhook support
- [ ] Middleware / global handlers
- [ ] Non-command message responses

---

## 🧑‍💻 Contributing

Pull requests and feature suggestions are welcome! This project aims to be developer-friendly.

---

## ⚠️ Warning

This library is under development. Use with caution in production environments.