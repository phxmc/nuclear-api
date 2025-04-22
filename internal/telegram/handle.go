package telegram

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/utils"
	"github.com/orewaee/nuclear-api/internal/validator"
)

func (bot *Bot) handle(update tgbotapi.Update, chatState domain.ChatState) {
	ctx := context.TODO()
	chatId := update.Message.Chat.ID

	switch chatState {
	case domain.StateEnterEmail:
		email := update.Message.Text
		ok, err := regexp.MatchString(validator.EmailRegexp, email)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		if !ok {
			message := tgbotapi.NewMessage(chatId, "Введите почту")
			if _, err := bot.api.Send(message); err != nil {
				bot.log.Error().Err(err).Send()
				return
			}
			return
		}

		exists, err := bot.accountApi.AccountExistsByEmail(ctx, email)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		if !exists {
			message := tgbotapi.NewMessage(chatId, "У тебя нет аккаунта")
			if _, err := bot.api.Send(message); err != nil {
				bot.log.Error().Err(err).Send()
				return
			}
			return
		}

		code := utils.MustNewCode()
		err = bot.telegramApi.AddTempTelegram(ctx, chatId, code, email, time.Minute*5)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		go func() {
			err := bot.emailApi.Send(ctx, email, "Ваш код - "+code, "Подключите Telegram к своему Nuclear аккаунту")
			if err != nil {
				bot.log.Error().Err(err).Send()
			}
		}()

		message := tgbotapi.NewMessage(chatId, fmt.Sprintf("Мы отправили на почту %s код подтверждения", email))
		if _, err := bot.api.Send(message); err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		err = bot.telegramApi.SetChatState(ctx, chatId, domain.StateEnterCode, time.Minute*5)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		break
	case domain.StateEnterCode:
		code := update.Message.Text

		tempTelegram, err := bot.telegramApi.GetTempTelegram(ctx, chatId)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		account, err := bot.accountApi.GetAccountByEmail(ctx, tempTelegram.Email)
		if errors.Is(err, domain.ErrNoAccount) {
			message := tgbotapi.NewMessage(chatId, "У тебя нет аккаунта")
			if _, err := bot.api.Send(message); err != nil {
				bot.log.Error().Err(err).Send()
				return
			}
			return
		}

		if code != tempTelegram.Code {
			message := tgbotapi.NewMessage(chatId, "Неверный код")
			if _, err := bot.api.Send(message); err != nil {
				bot.log.Error().Err(err).Send()
				return
			}
			return
		}

		err = bot.accountApi.SetAccountTelegramId(ctx, account.Id, chatId)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		message := tgbotapi.NewMessage(chatId, "Аккаунт успешно привязан! Используйте /me, чтобы узнать информацию об аккаунте")
		if _, err := bot.api.Send(message); err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		err = bot.telegramApi.ResetChatState(ctx, chatId)
		if err != nil {
			bot.log.Error().Err(err).Send()
			return
		}

		// todo reset temp telegram

		break
	}
}
