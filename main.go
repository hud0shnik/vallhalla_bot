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
var botUrl string

// Структуры для работы с Telegram API

type telegramResponse struct {
	Result []update `json:"result"`
}

type update struct {
	UpdateId int     `json:"update_id"`
	Message  message `json:"message"`
}

type message struct {
	Chat    chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker sticker `json:"sticker"`
}

type sticker struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
}

type chat struct {
	ChatId int `json:"id"`
}

type sendMessage struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type sendSticker struct {
	ChatId     int    `json:"chat_id"`
	StickerUrl string `json:"sticker"`
}

// Структура респонса vall-halla-api
type infoResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Drinks  []drinkInfo `json:"result"`
}

// Структура коктейля
type drinkInfo struct {
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
func sendMsg(chatId int, text string) error {

	// Формирование сообщения
	buf, err := json.Marshal(sendMessage{
		ChatId:    chatId,
		Text:      text,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка сообщения
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("sendMessage error: %s", err)
		return err
	}

	return nil

}

// Функция отправки стикера
func sendStck(chatId int, stickerId string) error {

	// Формирование стикера
	buf, err := json.Marshal(sendSticker{
		ChatId:     chatId,
		StickerUrl: stickerId,
	})
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка стикера
	_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
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

	// Проверка статускода
	switch resp.StatusCode {
	case 200:
		// При хорошем статусе респонса продолжение выполнения кода
	case 404:
		sendMsg(chatId, "drinks not found")
		sendStck(chatId, "CAACAgIAAxkBAAIBx2PriuCsDDVv8tcdbqZ42v90M8WeAAIzAQAC5y5hCNndnbfZVPwxLgQ")
		return
	case 400:
		sendMsg(chatId, "bad request")
		return
	default:
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
	var response infoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("json.Unmarshal error: %s", err)
		sendMsg(chatId, "internal error")
		return
	}

	// Отправка коктейлей
	for _, drink := range response.Drinks {
		sendMsg(chatId, fmt.Sprintf(
			"<b><pre>%s</pre></b>\nIt's a <b>%s</b>, <b>%s</b> and <b>%s</b> drink coasting <b>$%d</b>\n"+
				"<b>Recipe</b> - %s\n<b>Shortcut</b> - <u>%s</u>\n\n<i>\"%s\"</i>",
			drink.Name, drink.Flavour, drink.Primary_Type, drink.Secondary_Type, drink.Price, drink.Recipe, drink.Shortcut, drink.Description))
	}

}

// Функция генерации и отправки ответа
func respond(update update) {

	// Проверка на сообщение
	if update.Message.Text == "" {
		sendStck(update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBI2PrfKjtI2x-jY1WAs5MjFRBm6JwAAInAAOldjoOspa6vsFKQhkuBA")
		return
	}

	// Проверка на манипуляции с БД
	if strings.ContainsAny(update.Message.Text, "`'%+$") {
		sendMsg(update.Message.Chat.ChatId, "Bad request")
		return
	}

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
}

// Функция получения апдейтов
func getUpdates(offset int) ([]update, error) {

	// Запрос для получения апдейтов
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse telegramResponse
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
	botUrl = "https://api.telegram.org/bot" + viper.GetString("token")
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
