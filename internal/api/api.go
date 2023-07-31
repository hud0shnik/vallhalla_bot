package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hud0shnik/VallHalla_bot/internal/telegram"
	"github.com/sirupsen/logrus"
)

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

// Функция отправки рецептов
func SearchDrinks(botUrl string, chatId int, parameters []string) {

	// Запрос для получения рецептов
	resp, err := http.Get("https://vall-halla-api.vercel.app/api/info?" + strings.Join(parameters[1:], "&"))
	if err != nil {
		logrus.Printf("http.Get error: %s", err)
		telegram.SendMsg(botUrl, chatId, "vall-halla-api error")
		return
	}
	defer resp.Body.Close()

	// Проверка статускода
	switch resp.StatusCode {
	case 200:
		// При хорошем статусе респонса продолжение выполнения кода
	case 404:
		telegram.SendMsg(botUrl, chatId, "drinks not found")
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIBx2PriuCsDDVv8tcdbqZ42v90M8WeAAIzAQAC5y5hCNndnbfZVPwxLgQ")
		return
	case 400:
		telegram.SendMsg(botUrl, chatId, "bad request")
		return
	default:
		telegram.SendMsg(botUrl, chatId, "internal error")
		return
	}

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Printf("ioutil.ReadAll error: %s", err)
		telegram.SendMsg(botUrl, chatId, "internal error")
		return
	}
	var response infoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Printf("json.Unmarshal error: %s", err)
		telegram.SendMsg(botUrl, chatId, "internal error")
		return
	}

	// Отправка коктейлей
	for _, drink := range response.Drinks {
		telegram.SendMsg(botUrl, chatId, fmt.Sprintf(
			"<b><pre>%s</pre></b>\nIt's a <b>%s</b>, <b>%s</b> and <b>%s</b> drink coasting <b>$%d</b>\n"+
				"<b>Recipe</b> - %s\n<b>Shortcut</b> - <u>%s</u>\n\n<i>\"%s\"</i>",
			drink.Name, drink.Flavour, drink.Primary_Type, drink.Secondary_Type, drink.Price, drink.Recipe, drink.Shortcut, drink.Description))
	}

}
