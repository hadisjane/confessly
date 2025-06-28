package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" binding:"required" db:"username"`
	Email     string    `json:"email" binding:"required" db:"email"`
	Password  string    `json:"password" binding:"required" db:"password"`
	Role      string    `json:"role" db:"role"`
	Banned    bool      `json:"banned" db:"banned"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserRegister struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}