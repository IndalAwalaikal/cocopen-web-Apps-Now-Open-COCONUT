// controllers/test.go
package controllers

import (
	"database/sql"
	"cocopen-backend/dto"
	"cocopen-backend/middleware"
	"cocopen-backend/models"
	"cocopen-backend/services"
	"cocopen-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetUserSoalHandler mengambil semua soal untuk user
func GetUserSoalHandler(db *sql.DB) http.HandlerFunc {
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

        pendaftar, err := services.GetLatestPendaftarByUserID(db, claims.IDUser)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa pendaftaran: "+err.Error())
            return
        }
        if pendaftar == nil {
            utils.Error(w, http.StatusForbidden, "Anda belum melakukan pendaftaran")
            return
        }

        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.Error(w, http.StatusNotFound, "Tes belum tersedia")
                return
            }
            utils.Error(w, http.StatusInternalServerError, "Gagal memuat konfigurasi tes")
            return
        }

        now := time.Now()
        if testConfig.WaktuMulai != nil && now.Before(*testConfig.WaktuMulai) {
            utils.Error(w, http.StatusForbidden, "Tes belum dimulai")
            return
        }
        if testConfig.WaktuSelesai != nil && now.After(*testConfig.WaktuSelesai) {
            utils.Error(w, http.StatusForbidden, "Tes telah ditutup")
            return
        }

        alreadyTaken, err := services.HasUserTakenTest(db, claims.IDUser, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa riwayat tes")
            return
        }
        if alreadyTaken {
            utils.Error(w, http.StatusForbidden, "Anda sudah pernah mengikuti tes")
            return
        }

        soals, err := services.GetAllSoal(db)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal mengambil soal: "+err.Error())
            return
        }

        var response []dto.SoalResponse
        for _, s := range soals {
            response = append(response, dto.SoalResponse{
                IDSoal:     s.IDSoal,
                Nomor:      s.Nomor,
                Pertanyaan: s.Pertanyaan,
                PilihanA:   s.PilihanA,
                PilihanB:   s.PilihanB,
                PilihanC:   s.PilihanC,
                PilihanD:   s.PilihanD,
            })
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "soal":         response,
            "durasi_menit": testConfig.DurasiMenit,
            "judul":        testConfig.Judul,
            "deskripsi":    testConfig.Deskripsi,
        })
    }
}


func SubmitJawabanHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya metode POST yang diizinkan")
            return
        }

        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok {
            utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
            return
        }

        var req dto.SubmitJawabanRequest
        if err := utils.ParseAndValidate(r, &req); err != nil {
            utils.Error(w, http.StatusBadRequest, err.Error())
            return
        }

        pendaftar, err := services.GetLatestPendaftarByUserID(db, claims.IDUser)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa pendaftaran")
            return
        }
        if pendaftar == nil {
            utils.Error(w, http.StatusForbidden, "Anda belum mendaftar")
            return
        }
        pendaftarID := pendaftar.IDPendaftar

        // Ambil konfigurasi test
        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            utils.Error(w, http.StatusNotFound, "Tes tidak tersedia")
            return
        }

        // Cek apakah sudah pernah ujian
        alreadyTaken, err := services.HasUserTakenTest(db, claims.IDUser, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memeriksa riwayat tes")
            return
        }
        if alreadyTaken {
            utils.Error(w, http.StatusForbidden, "Anda sudah pernah mengikuti tes")
            return
        }

        // Ambil semua soal untuk penilaian
        soals, err := services.GetAllSoal(db)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal memuat soal untuk penilaian")
            return
        }

        // Mapping jawaban benar
        jawabanBenarMap := make(map[int]string)
        for _, s := range soals {
            jawabanBenarMap[s.IDSoal] = s.JawabanBenar
        }

        // Validasi jawaban
        for idSoal := range req.Jawaban {
            if _, exists := jawabanBenarMap[idSoal]; !exists {
                utils.Error(w, http.StatusBadRequest, "Soal dengan ID "+strconv.Itoa(idSoal)+" tidak ditemukan")
                return
            }
        }

        // Hitung skor
        skorBenar := 0
        var jawabans []models.JawabanUser

        for idSoal, jawabUser := range req.Jawaban {
            isBenar := strings.ToUpper(jawabUser) == jawabanBenarMap[idSoal]
            if isBenar {
                skorBenar++
            }
            jawabans = append(jawabans, models.JawabanUser{
                IDSoal:      idSoal,
                JawabanUser: strings.ToUpper(jawabUser),
                IsBenar:     isBenar,
            })
        }

        totalSoal := len(soals)
        skorSalah := totalSoal - skorBenar
        nilai := float64(skorBenar) / float64(totalSoal) * 100

        waktuSekarang := time.Now()

        hasilBaru := models.HasilTest{
            UserID:         claims.IDUser,
            PendaftarID:    pendaftarID, // ✅ gunakan dari pendaftar.IDPendaftar
            IDTest:         testConfig.IDTest,
            WaktuMulai:     waktuSekarang,
            WaktuSelesai:   &waktuSekarang,
            SkorBenar:      skorBenar,
            SkorSalah:      skorSalah,
            Nilai:          nilai,
        }

        err = services.CreateHasilTest(db, hasilBaru)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal membuat hasil tes: "+err.Error())
            return
        }

        idHasil, err := services.GetHasilIDByUserAndTest(db, claims.IDUser, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal ambil ID hasil tes")
            return
        }

        for _, j := range jawabans {
            j.IDHasil = idHasil
            err = services.CreateJawabanUser(db, j)
            if err != nil {
                utils.Error(w, http.StatusInternalServerError, "Gagal simpan jawaban: "+err.Error())
                return
            }
        }

        utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
            "success":    true,
            "message":    "Jawaban berhasil dikirim dan dinilai",
            "skor_benar": skorBenar,
            "skor_salah": skorSalah,
            "nilai":      nilai,
        })
    }
}

// CreateSoalHandler membuat soal baru (admin)
func CreateSoalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya POST yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		var req dto.SoalCreateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		soal := models.SoalTest{
			Nomor:        req.Nomor,
			Pertanyaan:   req.Pertanyaan,
			PilihanA:     req.PilihanA,
			PilihanB:     req.PilihanB,
			PilihanC:     req.PilihanC,
			PilihanD:     req.PilihanD,
			JawabanBenar: req.JawabanBenar,
		}

		err := services.CreateSoal(db, soal)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal membuat soal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"message": "Soal berhasil dibuat",
		})
	}
}

// UpdateSoalHandler memperbarui soal (admin)
func UpdateSoalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya PUT yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib")
			return
		}

		idSoal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		var req dto.SoalUpdateRequest
		if err := utils.ParseAndValidate(r, &req); err != nil {
			utils.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Ambil soal dari DB
		soal, err := services.GetSoalByID(db, idSoal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Soal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil soal")
			return
		}

		// Update field yang tidak nil
		if req.Nomor != nil {
			soal.Nomor = *req.Nomor
		}
		if req.Pertanyaan != nil {
			soal.Pertanyaan = *req.Pertanyaan
		}
		if req.PilihanA != nil {
			soal.PilihanA = *req.PilihanA
		}
		if req.PilihanB != nil {
			soal.PilihanB = *req.PilihanB
		}
		if req.PilihanC != nil {
			soal.PilihanC = *req.PilihanC
		}
		if req.PilihanD != nil {
			soal.PilihanD = *req.PilihanD
		}
		if req.JawabanBenar != nil {
			soal.JawabanBenar = *req.JawabanBenar
		}

		err = services.UpdateSoal(db, soal)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal update soal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Soal berhasil diperbarui",
		})
	}
}

// DeleteSoalHandler menghapus soal (admin)
func DeleteSoalHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya DELETE yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok || claims.Role != "admin" {
			utils.Error(w, http.StatusForbidden, "Akses ditolak")
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			utils.Error(w, http.StatusBadRequest, "Parameter id wajib")
			return
		}

		idSoal, err := strconv.Atoi(idStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "ID tidak valid")
			return
		}

		// Cek eksistensi soal
		_, err = services.GetSoalByID(db, idSoal)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Soal tidak ditemukan")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal cek soal")
			return
		}

		err = services.DeleteSoal(db, idSoal)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Gagal hapus soal: "+err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Soal berhasil dihapus",
		})
	}
}

// GetHasilTesUserHandler mendapatkan hasil tes user (user lihat hasil sendiri)
func GetHasilTesUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
			return
		}

		claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
		if !ok {
			utils.Error(w, http.StatusUnauthorized, "Akses ditolak")
			return
		}

		testConfig, err := services.GetTestConfig(db)
		if err != nil {
			utils.Error(w, http.StatusNotFound, "Konfigurasi tes tidak ditemukan")
			return
		}

		hasil, err := services.GetHasilByUserID(db, claims.IDUser, testConfig.IDTest)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.Error(w, http.StatusNotFound, "Anda belum mengikuti tes")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Gagal ambil hasil")
			return
		}

		var waktuSelesai *string
		if hasil.WaktuSelesai != nil {
			t := hasil.WaktuSelesai.Format("2006-01-02 15:04:05")
			waktuSelesai = &t
		}

		response := dto.HasilResponse{
			IDHasil:      hasil.IDHasil,
			UserID:       hasil.UserID,
			PendaftarID:  hasil.PendaftarID,
			SkorBenar:    hasil.SkorBenar,
			SkorSalah:    hasil.SkorSalah,
			Nilai:        hasil.Nilai,
			WaktuMulai:   hasil.WaktuMulai.Format("2006-01-02 15:04:05"),
			WaktuSelesai: waktuSelesai,
			DurasiMenit:  hasil.DurasiMenit,
		}

		utils.JSONResponse(w, http.StatusOK, response)
	}
}

// GetAllHasilTesHandler mengambil semua hasil tes (admin)
func GetAllHasilTesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            utils.Error(w, http.StatusMethodNotAllowed, "Hanya GET yang diizinkan")
            return
        }

        claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
        if !ok || claims.Role != "admin" {
            utils.Error(w, http.StatusForbidden, "Akses ditolak")
            return
        }

        testConfig, err := services.GetTestConfig(db)
        if err != nil {
            utils.Error(w, http.StatusNotFound, "Konfigurasi tes tidak ditemukan")
            return
        }

        rows, err := services.GetAllHasilTestByTestID(db, testConfig.IDTest)
        if err != nil {
            utils.Error(w, http.StatusInternalServerError, "Gagal ambil hasil tes: "+err.Error())
            return
        }
        defer rows.Close()

        var hasilList []dto.HasilResponse
        for rows.Next() {
            var h models.HasilTest
            var waktuSelesai sql.NullTime
            var durasiMenit sql.NullInt64

            err := rows.Scan(
                &h.IDHasil,
                &h.UserID,
                &h.PendaftarID,
                &h.SkorBenar,
                &h.SkorSalah,
                &h.Nilai,
                &h.WaktuMulai,
                &waktuSelesai,
                &durasiMenit,
            )
            if err != nil {
                utils.Error(w, http.StatusInternalServerError, "Gagal scan hasil: "+err.Error())
                return
            }

            var waktuSelesaiStr *string
            if waktuSelesai.Valid {
                t := waktuSelesai.Time.Format("2006-01-02 15:04:05")
                waktuSelesaiStr = &t
            }

            var durasi *int
            if durasiMenit.Valid {
                d := int(durasiMenit.Int64)
                durasi = &d
            }

            hasilList = append(hasilList, dto.HasilResponse{
                IDHasil:      h.IDHasil,
                UserID:       h.UserID,
                PendaftarID:  h.PendaftarID,
                SkorBenar:    h.SkorBenar,
                SkorSalah:    h.SkorSalah,
                Nilai:        h.Nilai,
                WaktuMulai:   h.WaktuMulai.Format("2006-01-02 15:04:05"),
                WaktuSelesai: waktuSelesaiStr,
                DurasiMenit:  durasi,
            })
        }

        utils.JSONResponse(w, http.StatusOK, hasilList)
    }
}