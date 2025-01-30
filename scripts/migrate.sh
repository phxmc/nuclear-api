#!/bin/bash
URL="postgres://root:kYJuSfL4FX7Mtcy2badzHpn9GmqUve6r@localhost:5442/nuclear?sslmode=disable"
migrate -path ./migrations -database "$URL" $1
