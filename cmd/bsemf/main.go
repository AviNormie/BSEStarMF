package main

import (
    "log"
    "net/http"

    "github.com/AviNormie/BSEStarMF/internal/api"
)

func main() {
    // Setup routes
    http.HandleFunc("/api/auth", api.AuthHandler)
    http.HandleFunc("/api/ucc/register", api.UCCRegistrationHandler)
    http.HandleFunc("/api/ucc/health", api.UCCHealthHandler)
    
    log.Println("BSE StAR MF API Server starting...")
    log.Println("Available endpoints:")
    log.Println("  POST /api/ucc/register - UCC Registration")
    log.Println("Server started at :8080")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
