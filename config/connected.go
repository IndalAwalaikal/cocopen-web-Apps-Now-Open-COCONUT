package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

func ConnectToDB() *sql.DB {
    LoadEnv()

    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        dbUser, dbPass, dbHost, dbPort, dbName,
    )

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Gagal membuka koneksi database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Gagal terhubung ke database: %v", err)
    }

    log.Println("Koneksi database berhasil")
    return db
}
