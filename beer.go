package main

import (
	"beer/config"
	"beer/daemon"
	telegrambot "beer/telegram_bot"
)

func main() {
	config.InitDb()
	go daemon.ReduceSobriety()
	telegrambot.InitBot(true)
}
