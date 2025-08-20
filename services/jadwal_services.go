// services/jadwal.go
package services

import (
	"database/sql"
	"cocopen-backend/models"
	"time"
)

func scanRowToJadwal(row *sql.Row) (*models.Jadwal, error) {
    var j models.Jadwal
    var pendaftarID sql.NullInt64
    var catatan, alasan, jenisJadwal sql.NullString
    var tanggalD sql.NullTime
    var jamMulaiD, jamSelesaiD sql.NullString

    err := row.Scan(
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
        &alasan,
        &tanggalD,
        &jamMulaiD,
        &jamSelesaiD,
        &jenisJadwal,
        &j.CreatedAt,
        &j.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    if pendaftarID.Valid {
        id := int(pendaftarID.Int64)
        j.PendaftarID = &id
    }
    if catatan.Valid {
        j.Catatan = &catatan.String
    }
    if alasan.Valid {
        j.AlasanPerubahan = &alasan.String
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
    if jenisJadwal.Valid {
        j.JenisJadwal = jenisJadwal.String
    } else {
        j.JenisJadwal = "pribadi"
    }

    return &j, nil
}

func scanJadwalRows(rows *sql.Rows) ([]models.Jadwal, error) {
	var jadwals []models.Jadwal

	for rows.Next() {
		var j models.Jadwal
		var pendaftarID sql.NullInt64
		var catatan, alasan, jenisJadwal sql.NullString
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
			&catatan,
			&j.PengajuanPerubahan,
			&alasan,
			&tanggalD,
			&jamMulaiD,
			&jamSelesaiD,
			&jenisJadwal,
			&j.CreatedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if pendaftarID.Valid {
			id := int(pendaftarID.Int64)
			j.PendaftarID = &id
		}
		if catatan.Valid {
			j.Catatan = &catatan.String
		}
		if alasan.Valid {
			j.AlasanPerubahan = &alasan.String
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
		if jenisJadwal.Valid {
			j.JenisJadwal = jenisJadwal.String
		} else {
			j.JenisJadwal = "pribadi"
		}

		jadwals = append(jadwals, j)
	}

	return jadwals, rows.Err()
}

func CreateJadwal(db *sql.DB, jadwal models.Jadwal) error {
	if jadwal.JenisJadwal == "" {
		jadwal.JenisJadwal = "pribadi"
	}

	query := `
		INSERT INTO jadwal (
			user_id, pendaftar_id, tanggal, jam_mulai, jam_selesai, 
			tempat, konfirmasi_jadwal, catatan, pengajuan_perubahan, 
			alasan_perubahan, jenis_jadwal
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		jadwal.UserID,
		jadwal.PendaftarID,
		jadwal.Tanggal,
		jadwal.JamMulai,
		jadwal.JamSelesai,
		jadwal.Tempat,
		jadwal.KonfirmasiJadwal,
		jadwal.Catatan,
		jadwal.PengajuanPerubahan,
		jadwal.AlasanPerubahan,
		jadwal.JenisJadwal,
	)
	return err
}

func GetAllJadwal(db *sql.DB) ([]models.Jadwal, error) {
	query := `
		SELECT 
			id_jadwal, user_id, pendaftar_id,
			tanggal, jam_mulai, jam_selesai,
			tempat, konfirmasi_jadwal,
			catatan, pengajuan_perubahan, alasan_perubahan,
			tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
			jenis_jadwal,
			created_at, updated_at
		FROM jadwal
		ORDER BY tanggal DESC, jam_mulai ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJadwalRows(rows)
}

func GetJadwalByID(db *sql.DB, id int) (*models.Jadwal, error) {
    query := `
        SELECT 
            id_jadwal, user_id, pendaftar_id,
            tanggal, jam_mulai, jam_selesai,
            tempat, konfirmasi_jadwal,
            catatan, pengajuan_perubahan, alasan_perubahan,
            tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
            jenis_jadwal,
            created_at, updated_at
        FROM jadwal
        WHERE id_jadwal = ?
    `
    row := db.QueryRow(query, id) 
    return scanRowToJadwal(row)  
}

func GetJadwalByUserID(db *sql.DB, userID int) ([]models.Jadwal, error) {
	query := `
		SELECT 
			id_jadwal, user_id, pendaftar_id,
			tanggal, jam_mulai, jam_selesai,
			tempat, konfirmasi_jadwal,
			catatan, pengajuan_perubahan, alasan_perubahan,
			tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
			jenis_jadwal,
			created_at, updated_at
		FROM jadwal
		WHERE user_id = ?
		ORDER BY tanggal DESC, jam_mulai ASC
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJadwalRows(rows)
}

// GetJadwalUmum mengambil jadwal yang bersifat umum
func GetJadwalUmum(db *sql.DB) ([]models.Jadwal, error) {
	query := `
		SELECT 
			id_jadwal, user_id, pendaftar_id,
			tanggal, jam_mulai, jam_selesai,
			tempat, konfirmasi_jadwal,
			catatan, pengajuan_perubahan, alasan_perubahan,
			tanggal_diajukan, jam_mulai_diajukan, jam_selesai_diajukan,
			jenis_jadwal,
			created_at, updated_at
		FROM jadwal
		WHERE jenis_jadwal = 'umum'
		ORDER BY tanggal DESC, jam_mulai ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJadwalRows(rows)
}

func UpdateJadwal(db *sql.DB, jadwal *models.Jadwal) error {
	query := `
		UPDATE jadwal SET
			pendaftar_id = ?,
			tanggal = ?,
			jam_mulai = ?,
			jam_selesai = ?,
			tempat = ?,
			konfirmasi_jadwal = ?,
			catatan = ?,
			jenis_jadwal = ?,
			updated_at = NOW()
		WHERE id_jadwal = ?
	`
	_, err := db.Exec(
		query,
		jadwal.PendaftarID,
		jadwal.Tanggal,
		jadwal.JamMulai,
		jadwal.JamSelesai,
		jadwal.Tempat,
		jadwal.KonfirmasiJadwal,
		jadwal.Catatan,
		jadwal.JenisJadwal,
		jadwal.IDJadwal,
	)
	return err
}

func UpdatePengajuanPerubahan(
	db *sql.DB,
	idJadwal int,
	pengajuan bool,
	alasan *string,
	tanggalD *time.Time,
	jamMulaiD *string,
	jamSelesaiD *string,
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
	_, err := db.Exec(
		query,
		pengajuan,
		alasan,
		tanggalD,
		jamMulaiD,
		jamSelesaiD,
		idJadwal,
	)
	return err
}

func DeleteJadwal(db *sql.DB, idJadwal int) error {
	query := `DELETE FROM jadwal WHERE id_jadwal = ?`
	_, err := db.Exec(query, idJadwal)
	return err
}