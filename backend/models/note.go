package models

type User struct {
	ID int `json:"id" gorm:"primaryKey"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Note struct {
	ID int `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	User string `json:"user"`
	Text string `json:"text"`
}

