package main

// Simple auth service for demonstration
func authenticate(token string) bool {
    return token != ""
}