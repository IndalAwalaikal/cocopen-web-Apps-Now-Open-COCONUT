// routes/routes.go
package routes

import (
	"cocopen-backend/controllers"
	"cocopen-backend/middleware"
	"database/sql"
	"net/http"
)

func Setup(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Static file server
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// === Auth Routes ===
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

	// === Protected Routes - User Role ===
	mux.Handle("/profile", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProfile(db)(w, r)
	}))

	mux.Handle("/profile/update", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateProfile(db)(w, r)
	}))

	mux.Handle("/pendaftar/create", middleware.Auth(middleware.Role("user")(func(w http.ResponseWriter, r *http.Request) {
		controllers.CreatePendaftar(db)(w, r)
	})))

	// 🔹 User: Lihat SEMUA jadwal (pribadi + umum) → tidak perlu tahu jenisnya
	mux.Handle("/jadwal/user", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUserJadwalHandler(db)(w, r)
	}))

	// 🔹 User: Ajukan perubahan jadwal
	mux.Handle("/jadwal/ajukan", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		controllers.AjukanPerubahanJadwalHandler(db)(w, r)
	}))

	// 🔹 User: Batalkan pengajuan perubahan
	mux.Handle("/jadwal/cancel-perubahan", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		controllers.CancelPengajuanPerubahanHandler(db)(w, r)
	}))

	// === Protected Routes - Admin Role ===
	// Pendaftar
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

	// 🔹 Jadwal - Admin
	mux.Handle("/jadwal/delete", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteJadwalHandler(db)(w, r)
	})))

	mux.Handle("/jadwal/update", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateJadwalHandler(db)(w, r)
	})))

	mux.Handle("/jadwal/create", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateJadwalHandler(db)(w, r)
	})))

	mux.Handle("/jadwal/all", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllJadwalHandler(db)(w, r)
	})))

	// GET /jadwal?id=123
	mux.Handle("/jadwal", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetJadwalByIDHandler(db)(w, r)
	})))

	return mux
}