package main

import (
	"context"
	"os"
)

func main() {
	app := NewApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		PrintError(err.Error())
		os.Exit(1)
	}
}
