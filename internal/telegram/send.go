package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Структура для отправки сообщения
type sendMessage struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// Структура для отправки стикера
type sendSticker struct {
	ChatId     int    `json:"chat_id"`
	StickerUrl string `json:"sticker"`
}

// Функция отправки сообщения
func SendMsg(botUrl string, chatId int, text string) error {

	// Формирование сообщения
	buf, err := json.Marshal(sendMessage{
		ChatId:    chatId,
		Text:      text,
		ParseMode: "HTML",
	})
	if err != nil {
		logrus.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка сообщения
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		logrus.Printf("sendMessage error: %s", err)
		return err
	}

	return nil

}

// Функция отправки стикера
func SendStck(botUrl string, chatId int, stickerId string) error {

	// Формирование стикера
	buf, err := json.Marshal(sendSticker{
		ChatId:     chatId,
		StickerUrl: stickerId,
	})
	if err != nil {
		logrus.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка стикера
	_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		logrus.Printf("sendSticker error: %s", err)
		return err
	}

	return nil

}
