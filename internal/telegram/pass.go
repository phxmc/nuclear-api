package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

// /pass [add|remove] [account_id] [duration]

func (bot *Bot) pass(update tgbotapi.Update) {
	ctx := context.TODO()
	chatId := update.Message.Chat.ID

	// fmt.Println(chatId, update.Message.IsCommand(), update.Message.Command(), update.Message.CommandArguments())

	if update.Message == nil || !strings.HasPrefix(update.Message.Text, "/pass ") {
		return
	}

	rawArgs := strings.TrimPrefix(update.Message.Text, "/pass")
	rawArgs = strings.TrimSpace(rawArgs)

	args := strings.Split(rawArgs, " ")
	if len(args) != 3 && len(args) != 2 {
		return
	}

	account, err := bot.accountApi.GetAccountByTelegramId(ctx, chatId)
	if errors.Is(err, domain.ErrNoAccount) {
		bot.telegramApi.SendMessage(ctx, chatId, "you don't have an account")
		return
	}

	if err != nil {
		bot.log.Error().Err(err).Send()
		return
	}

	switch args[0] {
	case "add":
		fmt.Println("add")
		bot.addPass(ctx, account, args, update)
	case "remove":
		bot.removePass(ctx, account, args, update)
	}

}

func (bot *Bot) addPass(ctx context.Context, account *domain.Account, args []string, update tgbotapi.Update) {
	if !domain.HasPerm(account.Perms, domain.PermAddPass) &&
		!domain.HasPerm(account.Perms, domain.PermSuper) {
		fmt.Println("no permission")
		return
	}

	accountId := args[1]
	duration, err := time.ParseDuration(args[2])
	if err != nil {
		fmt.Println("cannot parse duration")
		return
	}

	from := time.Now()
	to := from.Add(duration)
	pass, err := bot.passApi.SetPass(ctx, accountId, &from, &to)
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.telegramApi.SendMessage(ctx, update.Message.Chat.ID, "added pass id: "+pass.Id)
}

func (bot *Bot) removePass(ctx context.Context, account *domain.Account, args []string, update tgbotapi.Update) {
	if !domain.HasPerm(account.Perms, domain.PermAddPass) &&
		!domain.HasPerm(account.Perms, domain.PermSuper) {
		return
	}

	accountId := args[1]
	fmt.Println(accountId)
}
