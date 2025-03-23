package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/nuclear-api/internal/validator"
	"regexp"
)

func (bot *Bot) handle(update tgbotapi.Update, chatState domain.ChatState) {
	switch chatState {
	case domain.StateEnterEmail:
		email := update.Message.Text
		ok, err := regexp.MatchString(validator.EmailRegexp, email)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		if !ok {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите почту")
			if _, err := bot.api.Send(message); err != nil {
				bot.log.Error().Err(err).Send()
				return
			}
			return
		}

		ctx := context.TODO()
		exists, err := bot.accountApi.AccountExistsByEmail(ctx, email)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		if !exists {
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "У тебя нет аккаунта")
			if _, err := bot.api.Send(message); err != nil {
				bot.log.Error().Err(err).Send()
				return
			}
			return
		}

		code := utils.MustNewCode()
		go func() {
			err := bot.emailApi.Send(ctx, email, "Ваш код - "+code, "Подключите Telegram к своему Nuclear аккаунту")
			if err != nil {
				bot.log.Error().Err(err).Send()
			}
		}()

		message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Мы отправили на почту %s код подтверждения", email))
		if _, err := bot.api.Send(message); err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		break
	case domain.StateEnterCode:
		break
	}
}
