package main

import (
	"github.com/zcubbs/oauth-showcase/client/cmd"
	"github.com/zcubbs/oauth-showcase/client/configs"
)

func init() {
	configs.Bootstrap()
}

func main() {
	// Print config
	configs.PrintConfig()

	// Start server
	cmd.Start()
}
