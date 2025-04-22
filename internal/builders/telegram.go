package builders

import (
	"github.com/orewaee/nuclear-api/internal/app/api"
	"github.com/orewaee/nuclear-api/internal/app/domain"
	"github.com/orewaee/nuclear-api/internal/app/repo"
	"github.com/orewaee/nuclear-api/internal/broker"
	"github.com/orewaee/nuclear-api/internal/services"
	"github.com/rs/zerolog"
)

type TelegramApiBuilder interface {
	Builder[api.TelegramApi]
	TelegramRepo(repo.TelegramReadWriter) TelegramApiBuilder
	TempTelegramRepo(repo.TempTelegramReadWriter) TelegramApiBuilder
	Broker(*broker.Broker[*domain.Message]) TelegramApiBuilder
	Log(*zerolog.Logger) TelegramApiBuilder
}

type telegramApiBuilder struct {
	telegramRepo     repo.TelegramReadWriter
	tempTelegramRepo repo.TempTelegramReadWriter
	broker           *broker.Broker[*domain.Message]
	log              *zerolog.Logger
}

func NewTelegramApiBuilder() TelegramApiBuilder {
	return &telegramApiBuilder{}
}

func (builder *telegramApiBuilder) Build() api.TelegramApi {
	return services.NewTelegramService(
		builder.telegramRepo,
		builder.tempTelegramRepo,
		builder.broker,
		builder.log,
	)
}

func (builder *telegramApiBuilder) TelegramRepo(telegramRepo repo.TelegramReadWriter) TelegramApiBuilder {
	builder.telegramRepo = telegramRepo
	return builder
}

func (builder *telegramApiBuilder) TempTelegramRepo(tempTelegramRepo repo.TempTelegramReadWriter) TelegramApiBuilder {
	builder.tempTelegramRepo = tempTelegramRepo
	return builder
}

func (builder *telegramApiBuilder) Broker(broker *broker.Broker[*domain.Message]) TelegramApiBuilder {
	builder.broker = broker
	return builder
}

func (builder *telegramApiBuilder) Log(log *zerolog.Logger) TelegramApiBuilder {
	builder.log = log
	return builder
}
