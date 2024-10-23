package main

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

// Middleware to add client IP to the context
func addClientIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from the request
		clientIP := r.RemoteAddr
		if ip := r.Header.Get("X-Real-IP"); ip != "" {
			clientIP = ip
		} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
			clientIP = ip
		}

		// Add IP to the request context

		ctx := context.WithValue(r.Context(), "client_ip", clientIP)

		newUUID, _ := uuid.NewUUID()

		ctx = context.WithValue(ctx, "uuid", newUUID)
		// Pass the request along with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
