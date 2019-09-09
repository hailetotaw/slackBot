package main

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func translateToAmharic(text []string) string {
	ctx := context.Background()
	// change the slice of string to string
	word := []string{strings.Join(text, "")}

	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the target language.
	target, err := language.Parse("am")
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	// Translates the text into Amharic.
	translations, err := client.Translate(ctx, word, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	return translations[0].Text

}
