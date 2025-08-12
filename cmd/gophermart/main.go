package main

import (
	"github.com/Nikolay961996/goferma/internal/server"
)

func main() {
	config := server.NewConfig()
	server.Run(config)
}
