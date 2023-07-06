package telegrambot

import (
	"beer/actions"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitBot(IsDebug bool) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI("key")
	if err != nil {
		fmt.Println(err)
	}
	bot.Debug = IsDebug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1

	updates := bot.GetUpdatesChan(u)
	var msg tgbotapi.MessageConfig
	for {
		if len(updates) != 0 {
			for update := range updates {
				if update.Message != nil {
					text := update.Message.Text
					groupID := update.Message.Chat.ID
					userID := update.Message.From.ID
					userName := update.Message.From.FirstName
					if isProfile(text) {
						sobr, sum := actions.GetData(groupID, userID)
						mes := fmt.Sprintf("%s, ваша трезвость - %d, выпито:\nпива - %.1fл\nводки - %.1fл\nвина - %.1fл", userName, sobr, sum.BeerSum, sum.VodkaSum, sum.WineSum)
						msg = tgbotapi.NewMessage(groupID, mes)
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
					} else if ok, drinkType := isDrink(text); ok {
						isFail, str := actions.Drink(groupID, userID, userName, drinkType)
						if isFail {
							var msgText string
							switch drinkType {
							case actions.BEER:
								msgText = fmt.Sprintf("%s, вам не удалось выпить пива по причине: \n%s", userName, str)
							case actions.VODKA:
								msgText = fmt.Sprintf("%s, вам не удалось выпить водки по причине: \n%s", userName, str)
							case actions.WINE:
								msgText = fmt.Sprintf("%s, вам не удалось выпить вина по причине: \n%s", userName, str)
							}
							msg = tgbotapi.NewMessage(groupID, msgText)
						} else {
							var msgText string
							switch drinkType {
							case actions.BEER:
								msgText = fmt.Sprintf("%s, пиво успешно выпито, трезвость уменьшена на 10", userName)
							case actions.VODKA:
								msgText = fmt.Sprintf("%s, водка успешно выпита, трезвость уменьшена на 50", userName)
							case actions.WINE:
								msgText = fmt.Sprintf("%s, вино успешно выпито, трезвость уменьшена на 30", userName)
							}
							msg = tgbotapi.NewMessage(groupID, msgText)
						}
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
					} else if isHelp(text) {
						mes := fmt.Sprintf("%s, список команд: \n/beer - выпить пива\n/vodka - выпить водки\n/wine - выпить вина\n/profile - получить профиль\n/top_drink - топ алкашей", userName)
						msg = tgbotapi.NewMessage(groupID, mes)
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
					} else if isTop(text) {
						mes := actions.GetTop(groupID)
						msg = tgbotapi.NewMessage(groupID, mes)
						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
					}
				} else if update.CallbackQuery != nil {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
					if _, err := bot.Request(callback); err != nil {
						panic(err)
					}

					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Неизвестный запрос")
					if _, err = bot.Send(msg); err != nil {
						panic(err)
					}
				}
			}

		}
	}
}
