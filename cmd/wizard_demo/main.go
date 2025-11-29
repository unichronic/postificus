package main

import (
	"fmt"
	"log"

	"postificus/internal/auth/wizard"
)

func main() {
	startURL := "https://dev.to/enter"
	successCriteria := "https://dev.to/"
	cookieName := "remember_user_token"

	fmt.Println("Starting Auth Wizard for Dev.to...")

	cookieValue, err := wizard.RunAuthFlow(startURL, successCriteria, cookieName)
	if err != nil {
		log.Fatalf("Auth Wizard failed: %v", err)
	}

	fmt.Println("\n--- SUCCESS ---")
	fmt.Printf("Extracted Cookie (%s): %s\n", cookieName, cookieValue)
	fmt.Println("You can now save this token to your database.")
}
