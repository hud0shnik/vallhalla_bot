package handler

import (
	"strings"

	"github.com/hud0shnik/VallHalla_bot/internal/api"
	"github.com/hud0shnik/VallHalla_bot/internal/send"
	"github.com/hud0shnik/VallHalla_bot/internal/telegram"
)

// Функция генерации и отправки ответа
func Respond(botUrl string, update telegram.Update) {

	// Проверка на сообщение
	if update.Message.Text == "" {
		send.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBI2PrfKjtI2x-jY1WAs5MjFRBm6JwAAInAAOldjoOspa6vsFKQhkuBA")
		return
	}

	// Проверка на манипуляции с БД
	if strings.ContainsAny(update.Message.Text, "`'%+$") {
		send.SendMsg(botUrl, update.Message.Chat.ChatId, "Bad request")
		return
	}

	// Разделение текста пользователя на слайс
	request := strings.Split(update.Message.Text, " ")

	// Обработчик команд
	switch request[0] {
	case "/search", "/info", "search", "s", "/s":
		api.SearchDrinks(botUrl, update.Message.Chat.ChatId, request)
	case "/help", "/start":
		send.SendMsg(botUrl, update.Message.Chat.ChatId, "Syntax:\n<b>/search alcoholic=no flavour=spicy</b> - <i>all non-alcoholic spicy drinks</i>\n<b>/search type=promo shortcut=3xT</b> - <i>all promo drinks with 3 Karmotrine</i>\n<b>/search name=piano</b> - <i>\"Piano Man\" and \"Piano Woman\" recieps</i>\n\n You can also use\n	<b>/search ice=yes&price=280&description=champaigne</b>")
		send.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBOmPrgHU_dc2p5aNX_s2tbo8MytiNAAKDAQAC5y5hCC7gW3lr-iVQLgQ")
	default:
		send.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIBRWPrgSjoO8gZfTKgA2N6vXGpo1fNAAK_AAPnLmEI82NgLSCbuiMuBA")
	}

}
