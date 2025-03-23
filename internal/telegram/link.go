package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"time"
)

func (bot *Bot) link(update tgbotapi.Update) {
	ctx := context.TODO()
	chatId := update.Message.Chat.ID

	err := bot.telegramApi.SetChatState(ctx, chatId, domain.StateEnterEmail, time.Minute)
	if err != nil {
		bot.log.Error().Err(err).Send()
		return
	}

	message := tgbotapi.NewMessage(chatId, "Введите свою почту")
	_, err = bot.api.Send(message)
	if err != nil {
		bot.log.Error().Err(err).Send()
		return
	}
}
