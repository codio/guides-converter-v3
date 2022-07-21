package main

import (
	"log"

	"github.com/codio/guides-converter-v3/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("guides-converter-v3 error: %s\n", err.Error())
	}
}
