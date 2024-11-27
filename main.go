package main

import (
	"go_rate_limit/middleware"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", middleware.RateLimiterMiddleware(http.HandlerFunc(handler)))
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
