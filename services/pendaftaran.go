// services/pendaftaran.go
package services

import (
    "database/sql"
	"cocopen-backend/dto"
	"cocopen-backend/models"
)

func CreatePendaftar(db *sql.DB, req dto.CreatePendaftarRequest) error {
    query := `
        INSERT INTO pendaftar (
            nama_lengkap, asal_kampus, prodi, semester, no_wa, domisili,
            alamat_sekarang, tinggal_dengan, alasan_masuk, pengetahuan_coconut,
            foto_path
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

    _, err := db.Exec(query,
        req.NamaLengkap,
        req.AsalKampus,
        req.Prodi,
        req.Semester,
        req.NoWA,
        req.Domisili,
        req.AlamatSekarang,
        req.TinggalDengan,
        req.AlasanMasuk,
        req.PengetahuanCoconut,
        req.FotoPath,
    )

    return err
}

func GetAllPendaftar(db *sql.DB) (*sql.Rows, error) {
    return db.Query(
        "SELECT id_pendaftar, nama_lengkap, asal_kampus, prodi, semester, no_wa, domisili, alamat_sekarang, tinggal_dengan, alasan_masuk, pengetahuan_coconut, foto_path, created_at, updated_at, status FROM pendaftar ORDER BY created_at DESC",
    )
}

func GetPendaftarByID(db *sql.DB, idPendaftar int) (models.Pendaftar, error) {
    var p models.Pendaftar
    err := db.QueryRow(`
        SELECT id_pendaftar, nama_lengkap, asal_kampus, prodi, semester, no_wa, domisili,
               alamat_sekarang, tinggal_dengan, alasan_masuk, pengetahuan_coconut, foto_path,
               created_at, updated_at, status
        FROM pendaftar
        WHERE id_pendaftar = ?
    `).Scan(
        &p.IDPendaftar, &p.NamaLengkap, &p.AsalKampus, &p.Prodi, &p.Semester, &p.NoWA,
        &p.Domisili, &p.AlamatSekarang, &p.TinggalDengan, &p.AlasanMasuk,
        &p.PengetahuanCoconut, &p.FotoPath, &p.CreatedAt, &p.UpdatedAt, &p.Status,
    )
    return p, err
}


func UpdatePendaftar(db *sql.DB, idPendaftar int, status string, fotoPath string) error {
	query := `
		UPDATE pendaftar
		SET status = ?, foto_path = ?, updated_at = NOW()
		WHERE id_pendaftar = ?
	`
	_, err := db.Exec(query, status, fotoPath, idPendaftar)
	return err
}



func DeletePendaftar(db *sql.DB, idPendaftar int) error {
    _, err := db.Exec("DELETE FROM pendaftar WHERE id_pendaftar = ?", idPendaftar)
    return err
}

