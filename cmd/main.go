package main

import (
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"os"
)

// init gets called before the main function
func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found create")
		os.Exit(1)
	}
}

func main() {
}
