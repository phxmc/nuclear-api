package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (bot *Bot) sendMessage(chatId int64, text string) {
	message := tgbotapi.NewMessage(chatId, text)
	message.ParseMode = tgbotapi.ModeHTML
	if _, err := bot.api.Send(message); err != nil {
		bot.log.Error().Err(err).Send()
	}
}
