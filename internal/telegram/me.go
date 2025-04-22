package telegram

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orewaee/nuclear-api/internal/app/domain"
)

func (bot *Bot) me(update tgbotapi.Update) {
	ctx := context.TODO()
	chatId := update.Message.Chat.ID
	account, err := bot.accountApi.GetAccountByTelegramId(ctx, chatId)
	if errors.Is(err, domain.ErrNoAccount) {
		bot.sendMessage(chatId, "У вас нет аккаунта.")
		return
	}

	text := fmt.Sprintf("id: <span class=\"tg-spoiler\"><code>%s</code></span>\n"+
		"email: <span class=\"tg-spoiler\"><code>%s</code></span>", account.Id, account.Email)

	bot.sendMessage(chatId, text)
}
