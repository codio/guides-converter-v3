package app

import (
	"log"

	"github.com/codio/guides-converter-v3/internal/assessments"
	"github.com/codio/guides-converter-v3/internal/content"
)

func Run() error {
	if err := assessments.Convert(); err != nil {
		log.Printf("error convert assessments: %s\n", err)
	}
	if err := content.Convert(); err != nil {
		log.Printf("error convert content: %s\n", err)
	}
	return nil
}
