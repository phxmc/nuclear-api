package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) info(update tgbotapi.Update) {
	bot.log.Debug().Str("command", "/info").Send()

	id := update.SentFrom().ID
	chatId := update.Message.Chat.ID
	text := fmt.Sprintf("id = <code>%d</code>\nchat_id = <code>%d</code>", id, chatId)
	message := tgbotapi.NewMessage(chatId, text)
	message.ParseMode = tgbotapi.ModeHTML

	if _, err := bot.api.Send(message); err != nil {
		bot.log.Error().Err(err).Send()
	}
}
