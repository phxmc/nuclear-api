package logger

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func NewZerolog() (*zerolog.Logger, error) {
	err := os.Mkdir("logs", os.ModePerm)

	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	now := time.Now().Format("02012006-150405MST")
	path := fmt.Sprintf("logs/%s.log", now)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	writer := zerolog.MultiLevelWriter(file, console)

	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &logger, nil
}
