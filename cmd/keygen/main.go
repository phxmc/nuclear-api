package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateAPIKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(bytes)
	return key[:length], nil
}

func main() {
	key, err := generateAPIKey(32)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("api key:", key)
}
