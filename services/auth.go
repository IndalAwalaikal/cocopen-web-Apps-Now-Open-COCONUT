package services

import (
    "database/sql"
    "time"
)

func Register(db *sql.DB, username, email, passwordHash, role string) (int, error) {
    res, err := db.Exec(
        "INSERT INTO users (username, full_name, email, password, profile_picture, role, is_verified) VALUES (?, ?, ?, ?, ?, ?, ?)",
        username, username, email, passwordHash, "default.jpg", role, false,
    )
    if err != nil {
        return 0, err
    }

    lastID, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(lastID), nil
}

func GetUserByUsername(db *sql.DB, username string) (int, string, string, string, string, error) {
    var id int
    var uname, email, hashed, role string

    err := db.QueryRow(`
        SELECT id_user, username, email, password, role
        FROM users
        WHERE username = ? AND is_verified = TRUE
    `, username).Scan(&id, &uname, &email, &hashed, &role)

    return id, uname, email, hashed, role, err
}

func GetUserByEmail(db *sql.DB, email string) (int, error) {
    var id int
    err := db.QueryRow("SELECT id_user FROM users WHERE email = ?", email).Scan(&id)
    return id, err
}

func DeleteUnverifiedUsersBefore(db *sql.DB, cutoff time.Time) error {
    _, err := db.Exec(`
        DELETE FROM users 
        WHERE is_verified = FALSE 
          AND created_at < ?
    `, cutoff)
    return err
}