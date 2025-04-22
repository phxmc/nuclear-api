package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

// /pass [add|remove] [user_id] [duration]

func (bot *Bot) pass(update tgbotapi.Update) {
	ctx := context.TODO()
	chatId := update.Message.Chat.ID

	// fmt.Println(chatId, update.Message.IsCommand(), update.Message.Command(), update.Message.CommandArguments())

	if update.Message == nil || !strings.HasPrefix(update.Message.Text, "/pass ") {
		return
	}

	args := strings.Split(strings.TrimPrefix(update.Message.Text, "/pass "), " ")
	fmt.Println(args)

	account, err := bot.accountApi.GetAccountByTelegramId(ctx, chatId)
	if errors.Is(err, domain.ErrNoAccount) {
		bot.telegramApi.SendMessage(ctx, chatId, "you don't have an account")
		return
	}

	if err != nil {
		bot.log.Error().Err(err).Send()
		return
	}

	domain.HasPerm(account.Perms, domain.PermAddPass)
}
