package entity

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
