package client

import (
	"arch-telegram-service/models"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

type Transformer struct{}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (t *Transformer) TransformUpdate(upd tgbotapi.Update) (*models.Update, error) {
	convertedUpd := models.Update{}
	convertedUpd.ID = upd.UpdateID
	convertedUpd.Message = &models.IncomingMessage{}

	if upd.CallbackQuery != nil {
		if err := t.transformCallbackData(upd.CallbackQuery.Data, convertedUpd.Message); err != nil {
			return nil, err
		}
		convertedUpd.Message.Username = upd.CallbackQuery.Message.From.UserName
		convertedUpd.Message.ChatID = upd.CallbackQuery.Message.Chat.ID
	}

	if upd.Message != nil {
		convertedUpd.Message.Text = upd.Message.Text
		convertedUpd.Message.Username = upd.Message.Chat.UserName
		convertedUpd.Message.ChatID = upd.Message.Chat.ID
		if upd.Message.Location != nil {
			convertedUpd.Message.Location = &models.Location{
				Longitude: fmt.Sprintf("%f", upd.Message.Location.Longitude),
				Latitude:  fmt.Sprintf("%f", upd.Message.Location.Latitude),
			}
		}
	}

	return &convertedUpd, nil
}

func (t *Transformer) transformCallbackData(data string, msg *models.IncomingMessage) error {
	res := strings.Split(data, " ")

	if res[0] == "/change_radius" {
		msg.Text = res[0]

		radius, err := strconv.Atoi(res[1])
		if err != nil {
			return fmt.Errorf("transform radius: %w", err)
		}
		msg.Radius = radius
	} else if res[0] == "/pagination" {
		msg.Text = res[0]

		cp, err := strconv.Atoi(res[1])
		if err != nil {
			return fmt.Errorf("transform current page: %w", err)
		}

		mp, err := strconv.Atoi(res[2])
		if err != nil {
			return fmt.Errorf("transform maxPages: %w", err)
		}

		offset, err := strconv.Atoi(res[3])
		if err != nil {
			return fmt.Errorf("transform offset: %w", err)
		}

		limit, err := strconv.Atoi(res[4])
		if err != nil {
			return fmt.Errorf("transform limit: %w", err)
		}

		radius, err := strconv.Atoi(res[5])
		if err != nil {
			return fmt.Errorf("transform radius: %w", err)
		}
		msg.Radius = radius

		msg.Location = &models.Location{
			Longitude: res[6],
			Latitude:  res[7],
		}

		msg.Paginator = &models.Paginator{
			MaxPages:    mp,
			CurrentPage: cp,
			Offset:      offset,
			Limit:       limit,
		}
	}

	return nil
}
