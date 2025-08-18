// controllers/jadwal.go
package controllers

import (
	"cocopen-backend/dto"
	"cocopen-backend/middleware"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
		if err := parseJadwalCreateRequest(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		tanggal, err := time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Format tanggal tidak valid")
			return
		}

		jadwal := models.Jadwal{
			UserID:             req.UserID,
			PendaftarID:        req.PendaftarID,
			Tanggal:            tanggal,
			JamMulai:           req.JamMulai,
			JamSelesai:         req.JamSelesai,
			Tempat:             req.Tempat,
			KonfirmasiJadwal:   "belum",
			Catatan:            req.Catatan,
			PengajuanPerubahan: false,
			AlasanPerubahan:    nil,
		}

		if err := services.CreateJadwal(db, jadwal); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal membuat jadwal")
			return
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Jadwal berhasil dibuat",
		})
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

		rows, err := services.GetAllJadwal(db) // Gunakan fungsi yang sudah diperbaiki
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal")
			return
		}
		defer rows.Close()

		var result []dto.JadwalAdminResponse

		for rows.Next() {
			var j models.Jadwal
			var pendaftarID sql.NullInt64
			var catatanStr, alasanStr sql.NullString
			var tanggalD sql.NullTime
			var jamMulaiD, jamSelesaiD sql.NullString

			err := rows.Scan(
				&j.IDJadwal,
				&j.UserID,
				&pendaftarID,
				&j.Tanggal,
				&j.JamMulai,
				&j.JamSelesai,
				&j.Tempat,
				&j.KonfirmasiJadwal,
				&catatanStr,
				&j.PengajuanPerubahan,
				&alasanStr,
				&tanggalD,
				&jamMulaiD,
				&jamSelesaiD,
				&j.CreatedAt,
				&j.UpdatedAt,
			)
			if err != nil {
				utils.Error(w, http.StatusInternalServerError, "Gagal membaca data jadwal")
				return
			}

			var pendaftarIDPtr *int
			if pendaftarID.Valid {
				id := int(pendaftarID.Int64)
				pendaftarIDPtr = &id
			}

			var catatanPtr *string
			if catatanStr.Valid {
				catatanPtr = &catatanStr.String
			}

			var alasanPtr *string
			if alasanStr.Valid {
				alasanPtr = &alasanStr.String
			}

			var tanggalDPtr *time.Time
			if tanggalD.Valid {
				tanggalDPtr = &tanggalD.Time
			}

			var jamMulaiDPtr *string
			if jamMulaiD.Valid {
				jamMulaiDPtr = &jamMulaiD.String
			}

			var jamSelesaiDPtr *string
			if jamSelesaiD.Valid {
				jamSelesaiDPtr = &jamSelesaiD.String
			}

			response := dto.JadwalAdminResponse{
				IDJadwal:               j.IDJadwal,
				UserID:                 j.UserID,
				PendaftarID:            pendaftarIDPtr,
				Tanggal:                j.Tanggal,
				JamMulai:               j.JamMulai,
				JamSelesai:             j.JamSelesai,
				Tempat:                 j.Tempat,
				KonfirmasiJadwal:       j.KonfirmasiJadwal,
				Catatan:                catatanPtr,
				PengajuanPerubahan:     j.PengajuanPerubahan,
				AlasanPerubahan:        alasanPtr,
				CreatedAt:              j.CreatedAt,
				UpdatedAt:              j.UpdatedAt,
				TanggalDiajukan:        tanggalDPtr,
				JamMulaiDiajukan:       jamMulaiDPtr,
				JamSelesaiDiajukan:     jamSelesaiDPtr,
			}

			result = append(result, response)
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
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal")
			return
		}

		var pendaftarIDPtr *int
		if jadwal.PendaftarID != nil {
			pendaftarIDPtr = jadwal.PendaftarID
		}

		var catatanPtr *string
		if jadwal.Catatan != nil {
			catatanPtr = jadwal.Catatan
		}

		var alasanPtr *string
		if jadwal.AlasanPerubahan != nil {
			alasanPtr = jadwal.AlasanPerubahan
		}

		response := dto.JadwalAdminResponse{
			IDJadwal:           jadwal.IDJadwal,
			UserID:             jadwal.UserID,
			PendaftarID:        pendaftarIDPtr,
			Tanggal:            jadwal.Tanggal,
			JamMulai:           jadwal.JamMulai,
			JamSelesai:         jadwal.JamSelesai,
			Tempat:             jadwal.Tempat,
			KonfirmasiJadwal:   jadwal.KonfirmasiJadwal,
			Catatan:            catatanPtr,
			PengajuanPerubahan: jadwal.PengajuanPerubahan,
			AlasanPerubahan:    alasanPtr,
			CreatedAt:          jadwal.CreatedAt,
			UpdatedAt:          jadwal.UpdatedAt,
		}

		utils.JSONResponse(w, http.StatusOK, response)
	}
}

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

		rows, err := services.GetJadwalByUserID(db, claims.IDUser)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil jadwal Anda")
			return
		}
		defer rows.Close()

		var result []dto.JadwalUserResponse

		for rows.Next() {
			var j models.Jadwal
			var catatanStr, alasanStr sql.NullString
			var tanggalD sql.NullTime
			var jamMulaiD, jamSelesaiD sql.NullString

			err := rows.Scan(
				&j.IDJadwal,
				&j.UserID,
				&j.PendaftarID,
				&j.Tanggal,
				&j.JamMulai,
				&j.JamSelesai,
				&j.Tempat,
				&j.KonfirmasiJadwal,
				&catatanStr,
				&j.PengajuanPerubahan,
				&alasanStr,
				&tanggalD,
				&jamMulaiD,
				&jamSelesaiD,
				&j.CreatedAt,
				&j.UpdatedAt,
			)
			if err != nil {
				utils.Error(w, http.StatusInternalServerError, "Gagal membaca data jadwal")
				return
			}

			var catatanPtr *string
			if catatanStr.Valid {
				catatanPtr = &catatanStr.String
			}

			var alasanPtr *string
			if alasanStr.Valid {
				alasanPtr = &alasanStr.String
			}

			var tanggalDPtr *time.Time
			if tanggalD.Valid {
				tanggalDPtr = &tanggalD.Time
			}

			var jamMulaiDPtr *string
			if jamMulaiD.Valid {
				jamMulaiDPtr = &jamMulaiD.String
			}

			var jamSelesaiDPtr *string
			if jamSelesaiD.Valid {
				jamSelesaiDPtr = &jamSelesaiD.String
			}

			result = append(result, dto.JadwalUserResponse{
				IDJadwal:           j.IDJadwal,
				Tanggal:            j.Tanggal,
				JamMulai:           j.JamMulai,
				JamSelesai:         j.JamSelesai,
				Tempat:             j.Tempat,
				KonfirmasiJadwal:   j.KonfirmasiJadwal,
				Catatan:            catatanPtr,
				PengajuanPerubahan: j.PengajuanPerubahan,
				AlasanPerubahan:    alasanPtr,
				TanggalDiajukan:    tanggalDPtr,
				JamMulaiDiajukan:   jamMulaiDPtr,
				JamSelesaiDiajukan: jamSelesaiDPtr,
			})
		}

		utils.JSONResponse(w, http.StatusOK, result)
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

		var req dto.JadwalUpdateRequest
		if err := utils.ParseJSON(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, "Data tidak valid")
			return
		}

		idJadwal := req.IDJadwal
		if idJadwal == 0 {
			utils.Error(w, http.StatusBadRequest, "ID jadwal tidak valid")
			return
		}

		existing, err := services.GetJadwalByID(db, idJadwal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal")
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

		if err := services.UpdateJadwal(db, existing); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal memperbarui jadwal")
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
        if err := utils.ParseJSON(r, &req); err != nil {
            utils.Error(w, http.StatusBadRequest, "Data tidak valid")
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
            utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal")
            return
        }

        if jadwal.UserID != claims.IDUser {
            utils.Error(w, http.StatusForbidden, "Bukan jadwal Anda")
            return
        }

        if jadwal.PengajuanPerubahan {
            utils.Error(w, http.StatusBadRequest, "Sudah ada pengajuan perubahan")
            return
        }

        err = services.UpdatePengajuanPerubahan(
            db,
            req.IDJadwal,
            true,
            &req.AlasanPerubahan,
            req.TanggalDiajukan,
            req.JamMulaiDiajukan,
            req.JamSelesaiDiajukan,
        )
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal ajukan perubahan jadwal")
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
			utils.Error(w, http.StatusInternalServerError, "Gagal mengambil data jadwal")
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

		err = services.UpdatePengajuanPerubahan(
			db,
			idJadwal,
			false,
			nil,
			nil, 
			nil, 
			nil, 
		)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal batalkan pengajuan")
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

        var req dto.DeleteJadwalRequest
		if err := utils.ParseJSON(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, "Body JSON tidak valid: "+err.Error())
			return
		}

        if req.IDJadwal == 0 {
            utils.Error(w, http.StatusBadRequest, "id_jadwal wajib diisi")
            return
        }

        _, err := services.GetJadwalByID(db, req.IDJadwal)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "Jadwal tidak ditemukan")
                return
            }
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa jadwal")
            return
        }

        if err := services.DeleteJadwal(db, req.IDJadwal); err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal menghapus jadwal")
            return
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success": true,
            "message": "Jadwal berhasil dihapus",
        })
    }
}



func parseJadwalCreateRequest(r *http.Request, req *dto.JadwalCreateRequest) error {
	if err := utils.ParseJSON(r, req); err != nil {
		return err
	}

	if req.UserID == 0 {
		return errors.New("user_id wajib diisi dan harus angka")
	}
	if req.Tanggal == "" {
		return errors.New("tanggal wajib diisi")
	}
	if req.JamMulai == "" {
		return errors.New("jam_mulai wajib diisi")
	}
	if req.JamSelesai == "" {
		return errors.New("jam_selesai wajib diisi")
	}
	if req.Tempat == "" {
		return errors.New("tempat wajib diisi")
	}

	return nil
}

// parseJadwalUpdateRequest: Parse JSON body untuk update
func parseJadwalUpdateRequest(r *http.Request, req *dto.JadwalUpdateRequest) error {
	return utils.ParseJSON(r, req)
}