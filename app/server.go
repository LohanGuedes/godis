package main

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/godis"
)

func main() {
	app := godis.Server{}

	if err := app.Start(); err != nil {
		fmt.Println("Godis failed to start: %s", err.Error())
	}
}
