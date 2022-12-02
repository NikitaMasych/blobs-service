package main

import (
	"os"

	"blobs/internal/cli"
)

// KV_VIPER_FILE=config.yaml go run main.go run service

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
