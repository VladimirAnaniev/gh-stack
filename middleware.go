package main

// Auth middleware that uses the authentication service
func authMiddleware(token string, next func()) {
    if authenticate(token) {
        next()
    } else {
        panic("unauthorized")
    }
}