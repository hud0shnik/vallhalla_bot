package main

import (
	"log"

	"github.com/hud0shnik/VallHalla_bot/internal/config"
	"github.com/hud0shnik/VallHalla_bot/internal/handler"
	"github.com/hud0shnik/VallHalla_bot/internal/telegram"

	"github.com/spf13/viper"
)

func main() {

	// Инициализация конфига
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("initConfig error: %s", err)
		return
	}

	// Url бота для отправки и приёма сообщений
	botUrl := "https://api.telegram.org/bot" + viper.GetString("token")
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
