package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/orewaee/typedenv"
)

func NewCode() string {
	alphabet := typedenv.String("CODE_ALPHABET")
	size := typedenv.Int("CODE_SIZE", 8)

	return gonanoid.MustGenerate(alphabet, size)
}
