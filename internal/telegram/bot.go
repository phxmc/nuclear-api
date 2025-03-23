package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/driving"
	"github.com/orewaee/nuclear-api/internal/broker"
	"github.com/rs/zerolog"
)

type Bot struct {
	broker      *broker.Broker[string]
	api         *tgbotapi.BotAPI
	accountApi  api.AccountApi
	telegramApi api.TelegramApi
	emailApi    api.EmailApi
	log         *zerolog.Logger
}

func NewBot(
	accountApi api.AccountApi,
	telegramApi api.TelegramApi,
	emailApi api.EmailApi,
	log *zerolog.Logger) driving.Bot {
	return &Bot{
		broker:      broker.New[string](),
		accountApi:  accountApi,
		telegramApi: telegramApi,
		emailApi:    emailApi,
		log:         log,
	}
}

func (bot *Bot) SendMessage(ctx context.Context, markdown string) error {
	bot.broker.Publish(markdown)
	return nil
}

func (bot *Bot) Run(ctx context.Context, token string) error {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.api = botApi
	bot.api.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.api.GetUpdatesChan(updateConfig)

	messages := bot.broker.Subscribe()
	defer bot.broker.Unsubscribe(messages)
	go func() {
		for {
			message := <-messages
			fmt.Println(message)
		}
	}()

	for update := range updates {
		message := update.Message
		if message == nil {
			continue
		}

		chatState, ok := bot.telegramApi.GetChatState(ctx, message.Chat.ID)
		if ok {
			bot.handle(update, chatState)
		} else {
			var cmd func(tgbotapi.Update)
			switch message.Text {
			case "/info":
				cmd = bot.info
			case "/help":
				cmd = bot.help
			case "/link":
				cmd = bot.link
			default:
				cmd = bot.plain
			}

			cmd(update)
		}
	}

	return nil
}
