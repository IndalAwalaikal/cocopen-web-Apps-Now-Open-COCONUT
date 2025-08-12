// routes/routes.go
package routes

import (
	"database/sql"
	"net/http"

	"cocopen-backend/controllers"
    "cocopen-backend/middleware"
)

func Setup(db *sql.DB) http.Handler {
    mux := http.NewServeMux()

    //auth router

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

    // Protected - User role
    mux.Handle("/pendaftar/create", middleware.Auth(middleware.Role("user")(func(w http.ResponseWriter, r *http.Request) {
        controllers.CreatePendaftar(db)(w, r)
    })))

    // Protected - Admin role
    mux.Handle("/pendaftar/all", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
        controllers.GetAllPendaftar(db)(w, r)
    })))
    mux.Handle("/pendaftar/", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
        controllers.GetPendaftarByID(db)(w, r)
    })))
    mux.Handle("/pendaftar/update", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
        controllers.UpdatePendaftar(db)(w, r)
    })))
    mux.Handle("/pendaftar/delete", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
        controllers.DeletePendaftar(db)(w, r)
    })))

    return mux
}