// controllers/jadwal.go
package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/middleware"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"database/sql"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetUserJadwalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok {
			utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
			return
		}

		pribadi, err := services.GetJadwalByUserID(db, claims.IDUser)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil jadwal pribadi: "+err.Error())
			return
		}

		umum, err := services.GetJadwalUmum(db)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil jadwal umum: "+err.Error())
			return
		}

		allJadwal := append(pribadi, umum...)

		sort.Slice(allJadwal, func(i, j int) bool {
			if allJadwal[i].Tanggal == allJadwal[j].Tanggal {
				return allJadwal[i].JamMulai < allJadwal[j].JamMulai
			}
			return allJadwal[i].Tanggal.Before(allJadwal[j].Tanggal)
		})

		var result []dto.JadwalUserResponse
		for _, j := range allJadwal {
			result = append(result, dto.JadwalUserResponse{
				IDJadwal:           j.IDJadwal,
				Tanggal:            j.Tanggal,
				JamMulai:           j.JamMulai,
				JamSelesai:         j.JamSelesai,
				Tempat:             j.Tempat,
				KonfirmasiJadwal:   j.KonfirmasiJadwal,
				Catatan:            j.Catatan,
				PengajuanPerubahan: j.PengajuanPerubahan,
				AlasanPerubahan:    j.AlasanPerubahan,
				TanggalDiajukan:    j.TanggalDiajukan,
				JamMulaiDiajukan:   j.JamMulaiDiajukan,
				JamSelesaiDiajukan: j.JamSelesaiDiajukan,
			})
		}

		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func GetJadwalByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		idJadwal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID jadwal tidak valid")
			return
		}

		jadwal, err := services.GetJadwalByID(db, idJadwal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal: "+err.Error())
			return
		}

		response := dto.JadwalAdminResponse{
			IDJadwal:             jadwal.IDJadwal,
			UserID:               jadwal.UserID,
			PendaftarID:          jadwal.PendaftarID,
			Tanggal:              jadwal.Tanggal,
			JamMulai:             jadwal.JamMulai,
			JamSelesai:           jadwal.JamSelesai,
			Tempat:               jadwal.Tempat,
			KonfirmasiJadwal:     jadwal.KonfirmasiJadwal,
			Catatan:              jadwal.Catatan,
			PengajuanPerubahan:   jadwal.PengajuanPerubahan,
			AlasanPerubahan:      jadwal.AlasanPerubahan,
			TanggalDiajukan:      jadwal.TanggalDiajukan,
			JamMulaiDiajukan:     jadwal.JamMulaiDiajukan,
			JamSelesaiDiajukan:   jadwal.JamSelesaiDiajukan,
			JenisJadwal:          jadwal.JenisJadwal,
			CreatedAt:            jadwal.CreatedAt,
			UpdatedAt:            jadwal.UpdatedAt,
		}

		utils.JSONResponse(w, http.StatusOK, response)
	}
}

func GetAllJadwalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode GET yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		jadwals, err := services.GetAllJadwal(db)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal: "+err.Error())
			return
		}

		var result []dto.JadwalAdminResponse
		for _, j := range jadwals {
			result = append(result, dto.JadwalAdminResponse{
				IDJadwal:             j.IDJadwal,
				UserID:               j.UserID,
				PendaftarID:          j.PendaftarID,
				Tanggal:              j.Tanggal,
				JamMulai:             j.JamMulai,
				JamSelesai:           j.JamSelesai,
				Tempat:               j.Tempat,
				KonfirmasiJadwal:     j.KonfirmasiJadwal,
				Catatan:              j.Catatan,
				PengajuanPerubahan:   j.PengajuanPerubahan,
				AlasanPerubahan:      j.AlasanPerubahan,
				TanggalDiajukan:      j.TanggalDiajukan,
				JamMulaiDiajukan:     j.JamMulaiDiajukan,
				JamSelesaiDiajukan:   j.JamSelesaiDiajukan,
				JenisJadwal:          j.JenisJadwal,
				CreatedAt:            j.CreatedAt,
				UpdatedAt:            j.UpdatedAt,
			})
		}

		utils.JSONResponse(w, http.StatusOK, result)
	}
}

func CreateJadwalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		var req dto.JadwalCreateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		tanggal, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Format tanggal tidak valid")
			return
		}

		jamMulai, err := time.Parse("15:04:05", req.JamMulai)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Format jam_mulai tidak valid")
			return
		}

		jamSelesai, err := time.Parse("15:04:05", req.JamSelesai)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Format jam_selesai tidak valid")
			return
		}

		if !jamSelesai.After(jamMulai) {
			utils.Error(w, http.StatusBadRequest, "Jam selesai harus setelah jam mulai")
			return
		}

		jadwal := models.Jadwal{
			UserID:             claims.IDUser,
			PendaftarID:        req.PendaftarID,
			Tanggal:            tanggal,
			JamMulai:           req.JamMulai,
			JamSelesai:         req.JamSelesai,
			Tempat:             req.Tempat,
			KonfirmasiJadwal:   "belum",
			Catatan:            req.Catatan,
			PengajuanPerubahan: false,
			AlasanPerubahan:    nil,
			JenisJadwal:        "pribadi",
		}

		if req.JenisJadwal == "umum" {
			jadwal.JenisJadwal = "umum"
		}

		if err := services.CreateJadwal(db, jadwal); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal membuat jadwal: "+err.Error())
			return
		}

		if jadwal.JenisJadwal == "umum" {
    		utils.BroadcastNotifikasi(
        		"Jadwal Umum Baru",
        		"Jadwal umum telah ditambahkan di "+jadwal.Tempat+" pada "+jadwal.Tanggal.Format("2 Jan 2006"),
    		)
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Jadwal berhasil dibuat",
		})
	}
}

func UpdateJadwalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		idJadwal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID jadwal tidak valid")
			return
		}

		var req dto.JadwalUpdateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		existing, err := services.GetJadwalByID(db, idJadwal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal: "+err.Error())
			return
		}

		if req.PendaftarID != nil {
			existing.PendaftarID = req.PendaftarID
		}
		if req.Tanggal != nil {
			tanggal, err := time.Parse("2006-01-02", *req.Tanggal)
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "Format tanggal tidak valid")
				return
			}
			existing.Tanggal = tanggal
		}
		if req.JamMulai != nil {
			existing.JamMulai = *req.JamMulai
		}
		if req.JamSelesai != nil {
			existing.JamSelesai = *req.JamSelesai
		}
		if req.Tempat != nil {
			existing.Tempat = *req.Tempat
		}
		if req.KonfirmasiJadwal != nil {
			status := *req.KonfirmasiJadwal
			if status != "belum" && status != "dikonfirmasi" && status != "ditolak" {
				utils.Error(w, http.StatusBadRequest, "Status konfirmasi tidak valid")
				return
			}
			existing.KonfirmasiJadwal = status
		}
		if req.Catatan != nil {
			existing.Catatan = req.Catatan
		}
		if req.JenisJadwal != nil {
			if *req.JenisJadwal != "pribadi" && *req.JenisJadwal != "umum" {
				utils.Error(w, http.StatusBadRequest, "Jenis jadwal harus 'pribadi' atau 'umum'")
				return
			}
			existing.JenisJadwal = *req.JenisJadwal
		}

		if err := services.UpdateJadwal(db, existing); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal memperbarui jadwal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Jadwal berhasil diperbarui",
		})
	}
}

func AjukanPerubahanJadwalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "user" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		var req dto.JadwalAjukanPerubahanRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(strings.TrimSpace(req.AlasanPerubahan)) < 10 {
			utils.Error(w, http.StatusBadRequest, "Alasan perubahan minimal 10 karakter")
			return
		}

		jadwal, err := services.GetJadwalByID(db, req.IDJadwal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal: "+err.Error())
			return
		}

		if jadwal.UserID != claims.IDUser {
			utils.Error(w, http.StatusForbidden, "Bukan jadwal Anda")
			return
		}

		if jadwal.PengajuanPerubahan {
			utils.Error(w, http.StatusBadRequest, "Sudah ada pengajuan perubahan aktif")
			return
		}

		var tanggalD *time.Time
		if req.TanggalDiajukan != nil {
			t, err := time.Parse("2006-01-02", *req.TanggalDiajukan)
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "Format tanggal_diajukan tidak valid")
				return
			}
			tanggalD = &t
		}

		err = services.UpdatePengajuanPerubahan(
			db,
			req.IDJadwal,
			true,
			&req.AlasanPerubahan,
			tanggalD,
			req.JamMulaiDiajukan,
			req.JamSelesaiDiajukan,
		)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal ajukan perubahan jadwal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pengajuan perubahan jadwal berhasil dikirim",
		})
	}
}

func CancelPengajuanPerubahanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok {
			utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id_jadwal")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id_jadwal wajib diisi")
			return
		}

		idJadwal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID jadwal tidak valid")
			return
		}

		jadwal, err := services.GetJadwalByID(db, idJadwal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal: "+err.Error())
			return
		}

		if jadwal.UserID != claims.IDUser {
			utils.Error(w, http.StatusForbidden, "Anda tidak berhak mengakses jadwal ini")
			return
		}

		if !jadwal.PengajuanPerubahan {
			utils.Error(w, http.StatusBadRequest, "Tidak ada pengajuan perubahan untuk dibatalkan")
			return
		}

		err = services.UpdatePengajuanPerubahan(db, idJadwal, false, nil, nil, nil, nil)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal batalkan pengajuan: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Pengajuan perubahan berhasil dibatalkan",
		})
	}
}

func DeleteJadwalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib diisi")
			return
		}

		idJadwal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID jadwal tidak valid")
			return
		}

		_, err = services.GetJadwalByID(db, idJadwal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa jadwal: "+err.Error())
			return
		}

		if err := services.DeleteJadwal(db, idJadwal); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal menghapus jadwal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Jadwal berhasil dihapus",
		})
	}
}