package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/orewaee/typedenv"
)

func NewId() string {
	alphabet := typedenv.String("ID_ALPHABET")
	size := typedenv.Int("ID_SIZE", 8)

	return gonanoid.MustGenerate(alphabet, size)
}
