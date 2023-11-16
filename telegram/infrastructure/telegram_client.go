package client

import (
	"arch-telegram-service/models"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

type TelegramClient struct {
	Client *tgbotapi.BotAPI

	transformer   *Transformer
	buttonFactory *ButtonFactory
}

func NewTelegramClient(telegramToken string, paginationOffset int) *TelegramClient {
	client, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("I don't parse telegram token, err= %s", err)
	}

	transformer := NewTransformer()
	buttonFactory := NewButtonFactory(paginationOffset)

	return &TelegramClient{Client: client, transformer: transformer, buttonFactory: buttonFactory}
}

func (tc *TelegramClient) GetUpdates(offset, limit int) ([]models.Update, error) {
	uc := tgbotapi.UpdateConfig{
		Offset:  offset,
		Limit:   limit,
		Timeout: 60,
	}

	updates, err := tc.Client.GetUpdates(uc)
	if err != nil {
		return nil, fmt.Errorf("get updates from telegram: %w", err)
	}

	res := make([]models.Update, 0, len(updates))

	for _, update := range updates {
		convertedUpdate, err := tc.transformer.TransformUpdate(update)
		if err != nil {
			log.Printf("ошибка конвертации update в message, err: %s", err.Error())
			continue
		}

		res = append(res, *convertedUpdate)
	}

	return res, nil
}

func (tc *TelegramClient) SendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = tc.buttonFactory.CreateMainBtn()

	_, err := tc.Client.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message to telegram: %w", err)
	}

	return nil
}

func (tc *TelegramClient) SendRadiusMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = tc.buttonFactory.CreateRadiusInlineBtns()

	_, err := tc.Client.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message to telegram: %w", err)
	}

	return nil
}

func (tc *TelegramClient) SendPaginationBuildingsMessage(message models.OutgoingMessage) error {
	msg := tgbotapi.NewMessage(message.ChatID, message.Text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = tc.buttonFactory.CreatePaginationBuildingsInlineBtns(message)

	_, err := tc.Client.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send message to telegram: %w", err)
	}

	return nil
}
