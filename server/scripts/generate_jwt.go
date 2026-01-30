package main

import (
	"fmt"
	"server/internal/auth"
)

func main() {
	// Generate a token for user with ID 1 and role "admin"
	token, err := auth.GenerateToken(1, "admin", nil)
	if err != nil {
		fmt.Printf("Failed to generate token: %v ", err)
		return
	}

	fmt.Println("Generated JWT Token:")
	fmt.Println(token)
}
