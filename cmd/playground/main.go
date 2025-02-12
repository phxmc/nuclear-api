package main

import (
	"fmt"
	"github.com/orewaee/nuclear-api/internal/utils"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(utils.MustNewId())
	}
}
