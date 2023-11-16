package client

import (
	"arch-telegram-service/models"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type ButtonFactory struct {
	offset int
}

func NewButtonFactory(offset int) *ButtonFactory {
	return &ButtonFactory{offset: offset}
}

func (bf *ButtonFactory) CreateMainBtn() tgbotapi.ReplyKeyboardMarkup {
	mainBtn := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("🏠 Отправить локацию"),
			tgbotapi.NewKeyboardButton("🔘 Установить радиус поиска"),
		),
	)

	return mainBtn
}

// TODO много логики
func (bf *ButtonFactory) CreateRadiusInlineBtns() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	inlineBtn := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("500м", "/change_radius 500"),
		tgbotapi.NewInlineKeyboardButtonData("1км", "/change_radius 1000"),
		tgbotapi.NewInlineKeyboardButtonData("2км", "/change_radius 2000"),
		tgbotapi.NewInlineKeyboardButtonData("5км", "/change_radius 5000"),
		tgbotapi.NewInlineKeyboardButtonData("10км", "/change_radius 10000"),
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, inlineBtn)

	return keyboard
}

func (bf *ButtonFactory) CreatePaginationBuildingsInlineBtns(msg models.OutgoingMessage) tgbotapi.InlineKeyboardMarkup {
	if msg.Paginator == nil || msg.Location == nil {
		return tgbotapi.InlineKeyboardMarkup{}
	}
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	inlineBtn := make([]tgbotapi.InlineKeyboardButton, 0)

	currentPageNum := msg.Paginator.CurrentPage
	maxPages := msg.Paginator.MaxPages
	radius := msg.Radius
	longitude := msg.Location.Longitude
	latitude := msg.Location.Latitude

	res := make([]tgbotapi.InlineKeyboardButton, 0)
	if currentPageNum < maxPages && currentPageNum > 1 {
		res = []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"Назад",
				// собираем формат для предыдущей старницы с сообщениями
				fmt.Sprintf(
					"/pagination %v %v %v %v %v %v %v",
					currentPageNum-1,
					maxPages,
					(currentPageNum-1)*bf.offset-bf.offset,
					(currentPageNum-1)*bf.offset,
					radius,
					longitude,
					latitude,
				),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				"Далее",
				// собираем формат для следующей старницы с сообщениями
				fmt.Sprintf(
					"/pagination %v %v %v %v %v %v %v",
					currentPageNum+1,
					maxPages,
					currentPageNum*bf.offset,
					currentPageNum*bf.offset+bf.offset,
					radius,
					longitude,
					latitude,
				),
			),
		}
	} else if currentPageNum <= 1 && maxPages <= 1 {
		res = []tgbotapi.InlineKeyboardButton{}
	} else if currentPageNum == maxPages {
		res = []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"Назад",
				fmt.Sprintf(
					"/pagination %v %v %v %v %v %v %v",
					currentPageNum-1,
					maxPages,
					(currentPageNum-1)*bf.offset-bf.offset,
					(currentPageNum-1)*bf.offset,
					radius,
					longitude,
					latitude,
				),
			),
		}
	} else if currentPageNum == 1 && maxPages > 1 {
		res = []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				"Далее",
				fmt.Sprintf(
					"/pagination %v %v %v %v %v %v %v",
					currentPageNum+1,
					maxPages,
					currentPageNum*bf.offset,
					currentPageNum*bf.offset+bf.offset,
					radius,
					longitude,
					latitude,
				),
			),
		}
	}

	inlineBtn = append(inlineBtn, res...)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, inlineBtn)

	return keyboard
}
