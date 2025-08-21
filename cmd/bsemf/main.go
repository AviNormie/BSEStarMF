package main

import (
    "log"
    "net/http"

    "github.com/BrokingSapphire/BSEStarMF/internal/api"
)

func main() {
    http.HandleFunc("/api/auth", api.AuthHandler)
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
