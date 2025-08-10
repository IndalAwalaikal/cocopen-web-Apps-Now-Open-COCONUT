package main

import (
    "fmt"
    "net/http"
    "os"

    "cocopen-backend/config"
    "cocopen-backend/routes"
    "cocopen-backend/middleware"
)

func main() {
    config.LoadEnv()

    db := config.ConnectToDB()
    defer db.Close()

    mux := routes.Setup(db)

    handler := middleware.Recovery(middleware.Cors(mux))

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server berjalan di port %s\n", port)
    err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
    if err != nil {
        panic(fmt.Sprintf("Error menjalankan server: %v", err))
    }
}
