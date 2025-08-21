package main

import (
    "log"
    "net/http"

    "github.com/BrokingSapphire/BSEStarMF/internal/api"
)

func main() {
    // Setup routes
    http.HandleFunc("/api/auth", api.AuthHandler)
    http.HandleFunc("/api/ucc/register", api.UCCRegistrationHandler)
    http.HandleFunc("/api/ucc/health", api.UCCHealthHandler)
    
    // Lumpsum order entry routes
    http.HandleFunc("/api/lumpsum/order", api.LumpsumOrderHandler)
    http.HandleFunc("/api/lumpsum/order/test", api.LumpsumOrderTestHandler) // Test endpoint
    http.HandleFunc("/api/lumpsum/health", api.LumpsumOrderHealthHandler)
    
    log.Println("BSE StAR MF API Server starting...")
    log.Println("Available endpoints:")
    log.Println("  POST /api/auth - Authentication")
    log.Println("  POST /api/ucc/register - UCC Registration")
    log.Println("  GET /api/ucc/health - UCC Health Check")
    log.Println("  POST /api/lumpsum/order - Lumpsum Order Entry")
    log.Println("  POST /api/lumpsum/order/test - Lumpsum Order Test (No BSE call)")
    log.Println("  GET /api/lumpsum/health - Lumpsum Health Check")
    log.Println("Server started at :8080")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}