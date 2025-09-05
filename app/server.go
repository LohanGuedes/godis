package main

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/godis"
)

func main() {
	app := godis.NewServer(godis.Config{
		Host:            "0.0.0.0",
		Port:            6379,
		CleanerInterval: 1 * time.Second,
	})

	if err := app.Start(); err != nil {
		fmt.Println("Godis failed to start:", err.Error())
	}
}
