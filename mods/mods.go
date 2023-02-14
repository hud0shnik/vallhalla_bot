package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Структура респонса
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

// Функция отправки рецептов
func SendDrinkInfo(botUrl string, update Update, parameters []string) {

	// Rest запрос для получения апдейтов
	resp, err := http.Get("https://vall-halla-api.vercel.app/api/info?" + strings.Join(parameters[1:], "&"))
	if err != nil {
		log.Printf("http.Get error: %s", err)
	}
	defer resp.Body.Close()

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll error: %s", err)
	}
	var response InfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("json.Unmarshal error: %s", err)

	}

	// Отправка коктейлей
	for _, drink := range response.Drinks {
		SendMsg(botUrl, update, fmt.Sprintf(
			"<pre>%s</pre>\nIt's a <b>%s</b>, <b>%s</b> and <b>%s</b> drink coasting <b>$%d</b>\n"+
				"<b>Recipe</b> - %s\n<b>Shortcut</b> - <u>%s</u>\n\n<i>\"%s\"</i>",
			drink.Name, drink.Flavour, drink.Primary_Type, drink.Secondary_Type, drink.Price, drink.Recipe, drink.Shortcut, drink.Description))

	}

}
