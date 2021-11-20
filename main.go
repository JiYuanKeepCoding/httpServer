package main

import (
	"httpServer/connector"
)

func main() {
	server := connector.NewSever()
	server.Run("8080")
}
