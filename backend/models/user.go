package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
