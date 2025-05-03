package users

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Message string `json:"message"`
}

type Error struct {
	Error string `json:"error"`
}
