package main

import (
	"log"
	"os"

	"github.com/hud0shnik/VallHalla_bot/internal/handler"
	"github.com/hud0shnik/VallHalla_bot/internal/telegram"
	"github.com/joho/godotenv"
)

func main() {

	// Загрузка переменных окружения
	godotenv.Load()

	// Url бота для отправки и приёма сообщений
	botUrl := "https://api.telegram.org/bot" + os.Getenv("TOKEN")
	offSet := 0

	for {

		// Получение апдейтов
		updates, err := telegram.GetUpdates(botUrl, offSet)
		if err != nil {
			log.Fatalf("getUpdates error: %s", err)
			return
		}

		// Обработка апдейтов
		for _, update := range updates {
			handler.Respond(botUrl, update)
			offSet = update.UpdateId + 1
		}

	}

}
