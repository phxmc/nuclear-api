package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/orewaee/typedenv"
)

func MustNewCode() string {
	alphabet := typedenv.String("CODE_ALPHABET")
	size := typedenv.Int("CODE_SIZE", 4)

	return gonanoid.MustGenerate(alphabet, size)
}
