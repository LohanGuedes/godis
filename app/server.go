package main

import (
	"github.com/codecrafters-io/redis-starter-go/internal/eventloop"
	"github.com/codecrafters-io/redis-starter-go/internal/godis"
)

func main() {
	app := godis.Server{}

	eventloop.Start()
}
