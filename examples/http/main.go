// Package main demonstrates using Souuup for HTTP request validation
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cachesdev/souuup/r"
	"github.com/cachesdev/souuup/u"
)

// User represents a user registration request
type UserRegistration struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Age             int    `json:"age"`
}

// Custom validation rules
func ValidEmail(fs u.FieldState[string]) error {
	email := fs.Value
	if !strings.Contains(email, "@") {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
}

func PasswordMatchRule(reg UserRegistration) u.Rule[string] {
	return func(fs u.FieldState[string]) error {
		if fs.Value != reg.ConfirmPassword {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	}
}

func StrongPasswordRule(fs u.FieldState[string]) error {
	password := fs.Value

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	if !strings.ContainsAny(password, "0123456789") {
		return fmt.Errorf("password must contain at least one digit")
	}

	return nil
}

// Handler for user registration
func registerHandler(w http.ResponseWriter, req *http.Request) {
	// Only allow POST requests
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var reg UserRegistration
	err := json.NewDecoder(req.Body).Decode(&reg)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create validation schema
	schema := u.Schema{
		"username": u.Field(reg.Username, r.NotZero, r.MinS(3), r.MaxS(20)),
		"email":    u.Field(reg.Email, r.NotZero, ValidEmail),
		"password": u.Field(reg.Password, r.NotZero, StrongPasswordRule, PasswordMatchRule(reg)),
		"age":      u.Field(reg.Age, r.MinN(18)),
	}

	// Create validator
	s := u.NewSouuup(schema)

	// Validate
	if err := s.Validate(); err != nil {
		// Return validation errors as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
		return
	}

	// If validation passes, process the registration
	// (in a real app, this would save the user to a database)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "User registered successfully",
	})
}

func main() {
	// Register HTTP handler
	http.HandleFunc("/register", registerHandler)

	// Start HTTP server
	fmt.Println("HTTP Validation Example")
	fmt.Println("======================")
	fmt.Println("Server listening on http://localhost:8080")
	fmt.Println("To test, send a POST request to http://localhost:8080/register with JSON body:")
	fmt.Println(`{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "Password123",
  "confirmPassword": "Password123",
  "age": 25
}`)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
