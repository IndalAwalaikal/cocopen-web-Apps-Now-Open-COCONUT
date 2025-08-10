// models/user.go
package models

import "time"

type User struct {
    IDUser    int       `json:"id_user"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
