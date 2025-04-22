package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/driving"
	"github.com/orewaee/nuclear-api/internal/broker"
	"github.com/rs/zerolog"
)

type Bot struct {
	broker      *broker.Broker[*domain.Message]
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
	messageBroker *broker.Broker[*domain.Message],
	log *zerolog.Logger) driving.Bot {
	return &Bot{
		broker:      messageBroker,
		accountApi:  accountApi,
		telegramApi: telegramApi,
		emailApi:    emailApi,
		log:         log,
	}
}

func (bot *Bot) listenTelegramApiMessages() {
	messages := bot.broker.Subscribe()
	defer bot.broker.Unsubscribe(messages)

	for message := range messages {
		tgMsg := tgbotapi.NewMessage(message.ChatId, message.Markdown)
		tgMsg.ParseMode = "Markdown"
		if _, err := bot.api.Send(tgMsg); err != nil {
			bot.log.Error().Err(err).Msg("failed to send message")
		}
	}
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

	go bot.listenTelegramApiMessages()

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
			case "/me":
				cmd = bot.me
			case "/pass":
				cmd = bot.pass
			default:
				cmd = bot.plain
			}

			cmd(update)
		}
	}

	return nil
}
