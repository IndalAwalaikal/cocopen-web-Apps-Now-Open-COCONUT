// routes/routes.go
package routes

import (
    "net/http"
    "database/sql"

    "cocopen-backend/controllers"
    "cocopen-backend/middleware"
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

    // 🔐 Protected routes - Harus login & role sesuai
    mux.HandleFunc("/admin/dashboard", 
        middleware.Auth(
            middleware.Role("admin")(
                controllers.AdminDashboard(),
            ),
        ),
    )

    mux.HandleFunc("/user/dashboard", 
        middleware.Auth(
            middleware.Role("user")(
                controllers.UserDashboard(),
            ),
        ),
    )

    return mux
}