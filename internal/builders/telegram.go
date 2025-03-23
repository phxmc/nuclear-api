package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type TelegramApiBuilder interface {
	Builder[api.TelegramApi]
	TelegramRepo(repo.TelegramReadWriter) TelegramApiBuilder
	Log(*zerolog.Logger) TelegramApiBuilder
}

type telegramApiBuilder struct {
	telegramRepo repo.TelegramReadWriter
	log          *zerolog.Logger
}

func NewTelegramApiBuilder() TelegramApiBuilder {
	return &telegramApiBuilder{}
}

func (builder *telegramApiBuilder) Build() api.TelegramApi {
	return services.NewTelegramService(
		builder.telegramRepo,
		builder.log,
	)
}

func (builder *telegramApiBuilder) TelegramRepo(telegramRepo repo.TelegramReadWriter) TelegramApiBuilder {
	builder.telegramRepo = telegramRepo
	return builder
}

func (builder *telegramApiBuilder) Log(log *zerolog.Logger) TelegramApiBuilder {
	builder.log = log
	return builder
}
