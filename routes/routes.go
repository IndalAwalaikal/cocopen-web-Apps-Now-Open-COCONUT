// routes/routes.go
package routes

import (
	"database/sql"
	"net/http"

	"cocopen-backend/controllers"
)

func Setup(db *sql.DB) http.Handler {
    mux := http.NewServeMux()

    // 🔓 Public routes - Tidak perlu login
    mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        controllers.Register(db)(w, r)
    })

    mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        controllers.Login(db)(w, r)
    })

    mux.HandleFunc("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
        controllers.ForgotPassword(db)(w, r)
    })

    mux.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
        controllers.ResetPassword(db)(w, r)
    })

    mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
        controllers.VerifyEmail(db)(w, r)
    })

    return mux
}