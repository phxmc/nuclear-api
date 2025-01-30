package main

import (
	"fmt"
	"time"
)

func main() {
	timestamp := time.Now().Add(time.Minute * 2)

	fmt.Printf("%s", timestamp)
}
