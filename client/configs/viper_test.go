package configs

import "testing"

func TestBootstrap(t *testing.T) {
	Bootstrap()
}

func TestPrintConfig(t *testing.T) {
	PrintConfig()
}

func TestGetConfig(t *testing.T) {
	Bootstrap()
	PrintConfig()
}
