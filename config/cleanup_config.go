package config

import (
	"cocopen-backend/services"
	"database/sql"
	"log"
	"time"
)

func StartCleanupJob(db *sql.DB) {
    log.Println("Background job: Membersihkan akun belum diverifikasi setiap 10 menit")

    ticker := time.NewTicker(10 * time.Minute)
    go func() {
    for range ticker.C {
        cutoff := time.Now().Add(-30 * time.Minute)
        err := services.DeleteUnverifiedUsersBefore(db, cutoff)
        if err != nil {
            log.Printf("Gagal membersihkan akun belum diverifikasi: %v", err)
        } else {
            log.Printf("Pembersihan selesai: akun dibuat sebelum %v telah dicek", cutoff)
        }
    }
}()
}