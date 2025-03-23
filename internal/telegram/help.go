package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (bot *Bot) help(update tgbotapi.Update) {
	text := `
	/help - показать доступные команды
	/link - привязать аккаут
	`

	message := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	message.ParseMode = tgbotapi.ModeMarkdown

	if _, err := bot.api.Send(message); err != nil {
		bot.log.Error().Err(err).Send()
	}
}
