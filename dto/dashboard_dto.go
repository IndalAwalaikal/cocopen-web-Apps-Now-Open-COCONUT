package dto

type DashboardResponse struct {
    Message string      `json:"message"`
    User    UserSummary `json:"user"`
    Access  string      `json:"access"`
    Data    interface{} `json:"data,omitempty"`
}

type UserSummary struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Role     string `json:"role"`
}