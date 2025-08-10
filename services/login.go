package services

import (
    "database/sql"
    "cocopen-backend/models"
)

func Login(db *sql.DB, username string) (*models.User, string, error) {
    var user models.User
    var hashed string

    err := db.QueryRow(`
        SELECT id_user, username, email, password, role, created_at, updated_at
        FROM users
        WHERE username = ? AND is_verified = TRUE
    `, username).Scan(
        &user.IDUser,
        &user.Username,
        &user.Email,
        &hashed,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, "", err
    }

    return &user, hashed, nil
}
