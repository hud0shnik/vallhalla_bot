package mods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// Структуры для работы с Telegram API

type TelegramResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat    Chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker Sticker `json:"sticker"`
}

type Sticker struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type SendMessage struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type SendSticker struct {
	ChatId     int    `json:"chat_id"`
	StickerUrl string `json:"sticker"`
}

// Функция отправки сообщения
func SendMsg(botUrl string, update Update, msg string) error {

	// Формирование сообщения
	botMessage := SendMessage{
		ChatId:    update.Message.Chat.ChatId,
		Text:      msg,
		ParseMode: "HTML",
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		SendErrorMessage(botUrl, update, 2)
		return err
	}

	// Отправка сообщения
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendMessage method error: ", err)
		SendErrorMessage(botUrl, update, 5)
		return err
	}
	return nil
}

// Функция отправки стикера
func SendStck(botUrl string, update Update, url string) error {

	// Формирование стикера
	botStickerMessage := SendSticker{
		ChatId:     update.Message.Chat.ChatId,
		StickerUrl: url,
	}
	buf, err := json.Marshal(botStickerMessage)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		SendErrorMessage(botUrl, update, 2)
		return err
	}
	// Отправка стикера
	_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendSticker method error: ", err)
		SendErrorMessage(botUrl, update, 3)
		return err
	}
	return nil
}

// Функция отправки сообщений об ошибках
func SendErrorMessage(botUrl string, update Update, errorCode int) {

	// Генерация текста ошибки по коду
	var result string
	switch errorCode {
	case 1:
		result = "Ошибка работы API"
	case 2:
		result = "Ошибка работы json.Marshal()"
	case 3:
		result = "Ошибка работы метода SendSticker"
	case 4:
		result = "Ошибка работы метода SendPhoto"
	case 5:
		result = "Ошибка работы метода SendMessage"
	case 6:
		result = "Ошибка работы stickers.json"
	default:
		result = "Неизвестная ошибка"
	}

	// Анонимное оповещение меня
	var updateDanya Update
	updateDanya.Message.Chat.ChatId = viper.GetInt("DanyaChatId")
	SendMsg(botUrl, updateDanya, "Дань, тут у одного из пользователей "+result+", надеюсь он скоро тебе о ней напишет.")

	// Вывод ошибки пользователю с просьбой связаться со мной для её устранения
	result += ", пожалуйста свяжитесь с моим создателем для устранения проблемы\n\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik"
	SendMsg(botUrl, update, result)
}
