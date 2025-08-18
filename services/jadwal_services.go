package services

import (
	"cocopen-backend/models"
	"database/sql"
)

func CreateJadwal(db *sql.DB, j models.Jadwal) error {
	query := `
		INSERT INTO jadwal (
			user_id,
			pendaftar_id,
			tanggal,
			jam_mulai,
			jam_selesai,
			tempat,
			konfirmasi_jadwal,
			catatan,
			pengajuan_perubahan,
			alasan_perubahan,
			tanggal_diajukan,
			jam_mulai_diajukan,
			jam_selesai_diajukan
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var pendaftarID interface{} = j.PendaftarID
	if j.PendaftarID == nil {
		pendaftarID = nil
	}

	var catatan interface{} = j.Catatan
	if j.Catatan == nil {
		catatan = nil
	}

	var alasan interface{} = j.AlasanPerubahan
	if j.AlasanPerubahan == nil {
		alasan = nil
	}

	var tanggalD interface{} = j.TanggalDiajukan
	if j.TanggalDiajukan == nil {
		tanggalD = nil
	}

	var jamMulaiD interface{} = j.JamMulaiDiajukan
	if j.JamMulaiDiajukan == nil {
		jamMulaiD = nil
	}

	var jamSelesaiD interface{} = j.JamSelesaiDiajukan
	if j.JamSelesaiDiajukan == nil {
		jamSelesaiD = nil
	}

	_, err := db.Exec(
		query,
		j.UserID,
		pendaftarID,
		j.Tanggal,
		j.JamMulai,
		j.JamSelesai,
		j.Tempat,
		j.KonfirmasiJadwal,
		catatan,
		j.PengajuanPerubahan,
		alasan,
		tanggalD,
		jamMulaiD,
		jamSelesaiD,
	)
	return err
}

func GetAllJadwal(db *sql.DB) (*sql.Rows, error) {
	return db.Query(`
		SELECT 
			id_jadwal, user_id, pendaftar_id, tanggal, jam_mulai, jam_selesai,
			tempat, konfirmasi_jadwal, catatan, pengajuan_perubahan, alasan_perubahan,
			tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
			created_at, updated_at
		FROM jadwal
		ORDER BY tanggal DESC, jam_mulai ASC
	`)
}

func GetJadwalByID(db *sql.DB, idJadwal int) (models.Jadwal, error) {
	var j models.Jadwal

	var pendaftarID sql.NullInt64
	var catatan, alasanPerubahan sql.NullString
	var tanggalD sql.NullTime
	var jamMulaiD, jamSelesaiD sql.NullString

	err := db.QueryRow(`
		SELECT 
			id_jadwal, user_id, pendaftar_id, tanggal, jam_mulai, jam_selesai,
			tempat, konfirmasi_jadwal, catatan, pengajuan_perubahan, alasan_perubahan,
			tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
			created_at, updated_at
		FROM jadwal
		WHERE id_jadwal = ?
	`, idJadwal).Scan(
		&j.IDJadwal,
		&j.UserID,
		&pendaftarID,
		&j.Tanggal,
		&j.JamMulai,
		&j.JamSelesai,
		&j.Tempat,
		&j.KonfirmasiJadwal,
		&catatan,
		&j.PengajuanPerubahan,
		&alasanPerubahan,
		&tanggalD,
		&jamMulaiD,
		&jamSelesaiD,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		return j, err
	}

	if pendaftarID.Valid {
		id := int(pendaftarID.Int64)
		j.PendaftarID = &id
	}
	if catatan.Valid {
		j.Catatan = &catatan.String
	}
	if alasanPerubahan.Valid {
		j.AlasanPerubahan = &alasanPerubahan.String
	}
	if tanggalD.Valid {
		j.TanggalDiajukan = &tanggalD.Time
	}
	if jamMulaiD.Valid {
		j.JamMulaiDiajukan = &jamMulaiD.String
	}
	if jamSelesaiD.Valid {
		j.JamSelesaiDiajukan = &jamSelesaiD.String
	}

	return j, nil
}

func GetJadwalByUserID(db *sql.DB, userID int) (*sql.Rows, error) {
	return db.Query(`
		SELECT 
			id_jadwal, user_id, pendaftar_id, tanggal, jam_mulai, jam_selesai,
			tempat, konfirmasi_jadwal, catatan, pengajuan_perubahan, alasan_perubahan,
			tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
			created_at, updated_at
		FROM jadwal
		WHERE user_id = ?
		ORDER BY tanggal DESC, jam_mulai ASC
	`, userID)
}

func UpdateJadwal(db *sql.DB, j models.Jadwal) error {
	query := `
		UPDATE jadwal SET
			pendaftar_id = ?,
			tanggal = ?,
			jam_mulai = ?,
			jam_selesai = ?,
			tempat = ?,
			konfirmasi_jadwal = ?,
			catatan = ?,
			tanggal_diajukan = ?,
			jam_mulai_diajukan = ?,
			jam_selesai_diajukan = ?,
			updated_at = NOW()
		WHERE id_jadwal = ?
	`

	var pendaftarID interface{} = j.PendaftarID
	if j.PendaftarID == nil {
		pendaftarID = nil
	}

	var catatan interface{} = j.Catatan
	if j.Catatan == nil {
		catatan = nil
	}

	var tanggalD interface{} = j.TanggalDiajukan
	if j.TanggalDiajukan == nil {
		tanggalD = nil
	}

	var jamMulaiD interface{} = j.JamMulaiDiajukan
	if j.JamMulaiDiajukan == nil {
		jamMulaiD = nil
	}

	var jamSelesaiD interface{} = j.JamSelesaiDiajukan
	if j.JamSelesaiDiajukan == nil {
		jamSelesaiD = nil
	}

	_, err := db.Exec(
		query,
		pendaftarID,
		j.Tanggal,
		j.JamMulai,
		j.JamSelesai,
		j.Tempat,
		j.KonfirmasiJadwal,
		catatan,
		tanggalD,
		jamMulaiD,
		jamSelesaiD,
		j.IDJadwal,
	)
	return err
}

func UpdatePengajuanPerubahan(
	db *sql.DB,
	idJadwal int,
	pengajuan bool,
	alasan *string,
	tanggalD, jamMulaiD, jamSelesaiD *string,
) error {
	query := `
		UPDATE jadwal SET
			pengajuan_perubahan = ?,
			alasan_perubahan = ?,
			tanggal_diajukan = ?,
			jam_mulai_diajukan = ?,
			jam_selesai_diajukan = ?,
			updated_at = NOW()
		WHERE id_jadwal = ?
	`

	var a interface{} = alasan
	if alasan == nil {
		a = nil
	}

	var tD interface{} = tanggalD
	if tanggalD == nil {
		tD = nil
	}

	var jmD interface{} = jamMulaiD
	if jamMulaiD == nil {
		jmD = nil
	}

	var jsD interface{} = jamSelesaiD
	if jamSelesaiD == nil {
		jsD = nil
	}

	_, err := db.Exec(query, pengajuan, a, tD, jmD, jsD, idJadwal)
	return err
}

func DeleteJadwal(db *sql.DB, idJadwal int) error {
	_, err := db.Exec("DELETE FROM jadwal WHERE id_jadwal = ?", idJadwal)
	return err
}
