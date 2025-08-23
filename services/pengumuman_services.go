package services

import (
	"database/sql"
	"cocopen-backend/models"
)

// CreatePengumuman membuat pengumuman baru
func CreatePengumuman(db *sql.DB, p models.Pengumuman) error {
	query := `INSERT INTO pengumuman (judul, isi) VALUES (?, ?)`
	_, err := db.Exec(query, p.Judul, p.Isi)
	return err
}

// GetAllPengumuman mengambil semua pengumuman (terbaru dulu)
func GetAllPengumuman(db *sql.DB) ([]models.Pengumuman, error) {
	query := `SELECT id_pengumuman, judul, isi, created_at, updated_at FROM pengumuman ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pengumumans []models.Pengumuman
	for rows.Next() {
		var p models.Pengumuman
		err := rows.Scan(&p.IDPengumuman, &p.Judul, &p.Isi, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pengumumans = append(pengumumans, p)
	}
	return pengumumans, rows.Err()
}

// GetPengumumanByID mengambil satu pengumuman
func GetPengumumanByID(db *sql.DB, id int) (*models.Pengumuman, error) {
	query := `SELECT id_pengumuman, judul, isi, created_at, updated_at FROM pengumuman WHERE id_pengumuman = ?`
	row := db.QueryRow(query, id)

	var p models.Pengumuman
	err := row.Scan(&p.IDPengumuman, &p.Judul, &p.Isi, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdatePengumuman memperbarui pengumuman
func UpdatePengumuman(db *sql.DB, p *models.Pengumuman) error {
	query := `UPDATE pengumuman SET judul = ?, isi = ?, updated_at = NOW() WHERE id_pengumuman = ?`
	_, err := db.Exec(query, p.Judul, p.Isi, p.IDPengumuman)
	return err
}

// DeletePengumuman menghapus pengumuman
func DeletePengumuman(db *sql.DB, id int) error {
	query := `DELETE FROM pengumuman WHERE id_pengumuman = ?`
	_, err := db.Exec(query, id)
	return err
}