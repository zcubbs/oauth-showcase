package main

import (
	"github.com/zcubbs/oauth-showcase/server/cmd"
	"github.com/zcubbs/oauth-showcase/server/configs"
)

func init() {
	configs.Bootstrap()
}

func main() {
	configs.PrintConfig()

	// Start server
	cmd.StartServer()
}
