package main

import "log"

// Simple logging service for demonstration
func logInfo(message string) {
    log.Printf("[INFO] %s", message)
}

func logError(message string) {
    log.Printf("[ERROR] %s", message)
}