package dto

type DashboardResponse struct {
    Message string      `json:"message"`
    User    UserSummary `json:"user"`
    Access  string      `json:"access"`
    Data    interface{} `json:"data,omitempty"`
}

type UserSummary struct {
    IDUser   int    `json:"id_user"`
    Username string `json:"username"`
    Role     string `json:"role"`
}