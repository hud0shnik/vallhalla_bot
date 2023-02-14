package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"vallHallaBot/mods"

	"github.com/spf13/viper"
)

func main() {

	// Инициализация конфига
	err := InitConfig()
	if err != nil {
		fmt.Println("Config error: ", err)
		return
	}

	// Url бота для отправки и приёма сообщений
	botUrl := "https://api.telegram.org/bot" + viper.GetString("token")
	offSet := 0

	// Цикл работы бота
	for {

		// Получение апдейтов
		updates, err := getUpdates(botUrl, offSet)
		if err != nil {
			fmt.Println("Something went wrong: ", err)
		}

		// Обработка апдейтов
		for _, update := range updates {
			respond(botUrl, update)
			offSet = update.UpdateId + 1
		}

		// Вывод в консоль для тестов
		// fmt.Println(updates)
	}
}

// Функция получения апдейтов
func getUpdates(botUrl string, offset int) ([]mods.Update, error) {

	// Rest запрос для получения апдейтов
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
	var restResponse mods.TelegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

// Функция генерации и отправки ответа
func respond(botUrl string, update mods.Update) error {

	// Обработчик команд
	if update.Message.Text != "" {

		request := append(strings.Split(update.Message.Text, " "), "", "")

		switch request[0] {
		case "/search":
			mods.SendDrinkInfo(botUrl, update, request)
			return nil
		case "/help", "/start":
			mods.SendMsg(botUrl, update, "Syntax:\n<b>/search alcoholic=no flavour=spicy</b> - <i>all non-alcoholic spicy drinks</i>\n<b>/search type=promo shortcut=3xT</b> - <i>all promo drinks with 3 Karmotrine</i>\n<b>/search name=piano</b> - <i>\"Piano Man\" and \"Piano Woman\" recieps</i>\n\n You can also use\n	<b>/search ice=yes&price=280&description=champaigne</b>")
			mods.SendStck(botUrl, update, "CAACAgIAAxkBAAIBOmPrgHU_dc2p5aNX_s2tbo8MytiNAAKDAQAC5y5hCC7gW3lr-iVQLgQ")
			return nil
		}

		// Дефолтный ответ
		mods.SendStck(botUrl, update, "CAACAgIAAxkBAAIBRWPrgSjoO8gZfTKgA2N6vXGpo1fNAAK_AAPnLmEI82NgLSCbuiMuBA")
		return nil

	} else {

		// Если пользователь отправил не сообщение
		mods.SendStck(botUrl, update, "CAACAgIAAxkBAAIBI2PrfKjtI2x-jY1WAs5MjFRBm6JwAAInAAOldjoOspa6vsFKQhkuBA")
		return nil

	}
}

// Функция инициализации конфига (всех токенов)
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
