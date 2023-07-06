package main

import (
	"beer/config"
	"beer/daemon"
	telegrambot "beer/telegram_bot"
)

func main() {
	config.InitDb()
	daemon.FillCount()
	go daemon.ReduceSobriety()
	telegrambot.InitBot(true)
}
