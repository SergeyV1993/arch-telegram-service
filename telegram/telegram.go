package telegram

import (
	"arch-telegram-service/models"
	commandMessage "arch-telegram-service/telegram/text_factory"
	"context"
	"errors"
	"fmt"
)

type Client interface {
	GetUpdates(offset, limit int) ([]models.Update, error)
	SendMessage(chatId int64, text string) error
	SendRadiusMessage(chatId int64, text string) error
	SendPaginationBuildingsMessage(message models.OutgoingMessage) error
}

type ServiceConfig struct {
	Offset int // сдвиги сообщений в телеграмме
	Limit  int

	DefaultPaginationOffset int // сколько изначально опоказываем зданий по умолчанию
}

type Service struct {
	TelegramClient Client
	Config         ServiceConfig
}

func NewTelegramService(client Client, config ServiceConfig) *Service {
	return &Service{TelegramClient: client, Config: config}
}

func (ts *Service) SendMessage(_ context.Context, message models.OutgoingMessage) error {
	switch message.Type {
	case models.ErrorMessageType:
		return ts.TelegramClient.SendMessage(message.ChatID, commandMessage.CreateErrorMsg())
	case models.NotFoundMessageType:
		return ts.TelegramClient.SendMessage(message.ChatID, commandMessage.NotFoundBuildingsMsg())
	case models.StartMessageType:
		return ts.TelegramClient.SendMessage(message.ChatID, commandMessage.CreateStartMsg(message.Username))
	case models.HelpMessageType:
		return ts.TelegramClient.SendMessage(message.ChatID, commandMessage.CreateHelpMsg())
	case models.RadiusVariantMessageType:
		return ts.TelegramClient.SendRadiusMessage(message.ChatID, commandMessage.CreateVariantRadiusMsg())
	case models.RadiusChangeMessageType:
		return ts.TelegramClient.SendMessage(message.ChatID, commandMessage.CreateSuccessResponseChangeRadiusMsg())
	case models.BuildingsWithPaginationMessageType:
		return ts.sendBuildingsWithPagination(message)
	case models.BuildingsMessageType:
		return ts.sendAllBuildings(message)
	default:
		return ts.TelegramClient.SendMessage(message.ChatID, commandMessage.CreateUnknownMsg())
	}
}

func (ts *Service) GetUpdatesFromTelegram() ([]models.Update, error) {
	updates, err := ts.TelegramClient.GetUpdates(ts.Config.Offset, ts.Config.Limit)
	if err != nil {
		return nil, fmt.Errorf("get events from telegram: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	ts.Config = ServiceConfig{
		Offset: updates[len(updates)-1].ID + 1,
		Limit:  ts.Config.Limit,
	}

	return updates, nil
}

func (ts *Service) sendAllBuildings(message models.OutgoingMessage) error {
	currentPage := 1
	buildingsLimit := ts.Config.DefaultPaginationOffset
	maxPages := len(message.Buildings) / buildingsLimit

	if len(message.Buildings) < buildingsLimit {
		buildingsLimit = len(message.Buildings)
	} else {
		maxPages++
	}

	res := models.OutgoingMessage{
		Text:   commandMessage.CreateBuildingsMessage(message.Buildings[:buildingsLimit]),
		ChatID: message.ChatID,
		Radius: message.Radius,
		Paginator: &models.Paginator{
			MaxPages:    maxPages,
			CurrentPage: currentPage,
		},
		Location: message.Location,
	}

	return ts.TelegramClient.SendPaginationBuildingsMessage(res)
}

func (ts *Service) sendBuildingsWithPagination(message models.OutgoingMessage) error {
	if message.Location == nil {
		return errors.New("empty locations")
	}
	if message.Paginator == nil {
		return errors.New("empty paginator")
	}

	message.Text = commandMessage.CreateBuildingsMessage(message.Buildings)

	return ts.TelegramClient.SendPaginationBuildingsMessage(message)
}
