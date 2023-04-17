package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Глобальный Url бота
var BotUrl string

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

// Структура респонса vall-halla-api
type InfoResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Drinks  []DrinkInfo `json:"result"`
}

// Структура коктейля
type DrinkInfo struct {
	Name           string `json:"name"`
	Price          int    `json:"price"`
	Flavour        string `json:"flavour"`
	Primary_Type   string `json:"primary_type"`
	Secondary_Type string `json:"secondary_type"`
	Recipe         string `json:"recipe"`
	Shortcut       string `json:"shortcut"`
	Description    string `json:"description"`
}

// Функция инициализации конфига (всех токенов)
func initConfig() error {

	// Путь и имя файла
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

// Функция отправки сообщения
func sendMsg(chatId int, msg string) error {

	// Формирование сообщения
	buf, err := json.Marshal(SendMessage{
		ChatId:    chatId,
		Text:      msg,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка сообщения
	_, err = http.Post(BotUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("sendMessage error: %s", err)
		return err
	}

	return nil
}

// Функция отправки стикера
func sendStck(chatId int, url string) error {

	// Формирование стикера
	buf, err := json.Marshal(SendSticker{
		ChatId:     chatId,
		StickerUrl: url,
	})
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка стикера
	_, err = http.Post(BotUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("sendSticker error: %s", err)
		return err
	}

	return nil
}

// Функция отправки рецептов
func searchDrinks(chatId int, parameters []string) {

	// Запрос для получения рецептов
	resp, err := http.Get("https://vall-halla-api.vercel.app/api/info?" + strings.Join(parameters[1:], "&"))
	if err != nil {
		log.Printf("http.Get error: %s", err)
		sendMsg(chatId, "vall-halla-api error")
		return
	}
	defer resp.Body.Close()

	// Проверка статускода респонса
	if resp.StatusCode != 200 {
		log.Printf("vall-halla-api error: %d", resp.StatusCode)
		sendMsg(chatId, "internal error")
		return
	}

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll error: %s", err)
		sendMsg(chatId, "internal error")
		return
	}
	var response InfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("json.Unmarshal error: %s", err)
		sendMsg(chatId, "internal error")
		return
	}

	// Проверка на респонс
	if len(response.Drinks) == 0 {
		sendMsg(chatId, "Drinks not found")
		sendStck(chatId, "CAACAgIAAxkBAAIBx2PriuCsDDVv8tcdbqZ42v90M8WeAAIzAQAC5y5hCNndnbfZVPwxLgQ")
	} else {

		// Отправка коктейлей
		for _, drink := range response.Drinks {
			sendMsg(chatId, fmt.Sprintf(
				"<b><pre>%s</pre><b>\nIt's a <b>%s</b>, <b>%s</b> and <b>%s</b> drink coasting <b>$%d</b>\n"+
					"<b>Recipe</b> - %s\n<b>Shortcut</b> - <u>%s</u>\n\n<i>\"%s\"</i>",
				drink.Name, drink.Flavour, drink.Primary_Type, drink.Secondary_Type, drink.Price, drink.Recipe, drink.Shortcut, drink.Description))
		}
	}
}

// Функция генерации и отправки ответа
func respond(update Update) {

	// Проверка на сообщение
	if update.Message.Text != "" {

		// Разделение текста пользователя на слайс
		request := strings.Split(update.Message.Text, " ")

		// Обработчик команд
		switch request[0] {
		case "/search", "/info", "search", "s", "/s":
			searchDrinks(update.Message.Chat.ChatId, request)
		case "/help", "/start":
			sendMsg(update.Message.Chat.ChatId, "Syntax:\n<b>/search alcoholic=no flavour=spicy</b> - <i>all non-alcoholic spicy drinks</i>\n<b>/search type=promo shortcut=3xT</b> - <i>all promo drinks with 3 Karmotrine</i>\n<b>/search name=piano</b> - <i>\"Piano Man\" and \"Piano Woman\" recieps</i>\n\n You can also use\n	<b>/search ice=yes&price=280&description=champaigne</b>")
			sendStck(update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBOmPrgHU_dc2p5aNX_s2tbo8MytiNAAKDAQAC5y5hCC7gW3lr-iVQLgQ")
		default:
			sendStck(update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBRWPrgSjoO8gZfTKgA2N6vXGpo1fNAAK_AAPnLmEI82NgLSCbuiMuBA")
		}

	} else {

		// Если пользователь отправил не сообщение
		sendStck(update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBI2PrfKjtI2x-jY1WAs5MjFRBm6JwAAInAAOldjoOspa6vsFKQhkuBA")

	}
}

// Функция получения апдейтов
func getUpdates(offset int) ([]Update, error) {

	// Запрос для получения апдейтов
	resp, err := http.Get(BotUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse TelegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func main() {

	// Инициализация конфига
	err := initConfig()
	if err != nil {
		log.Fatalf("initConfig error: %s", err)
		return
	}

	// Url бота для отправки и приёма сообщений
	BotUrl = "https://api.telegram.org/bot" + viper.GetString("token")
	offSet := 0

	// Цикл работы бота
	for {

		// Получение апдейтов
		updates, err := getUpdates(offSet)
		if err != nil {
			log.Fatalf("getUpdates error: %s", err)
			return
		}

		// Обработка апдейтов
		for _, update := range updates {
			respond(update)
			offSet = update.UpdateId + 1
		}

	}
}
