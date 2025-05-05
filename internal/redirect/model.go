package redirect

import (
	"github.com/google/uuid"
	"time"
)

type Click struct {
	ID        uuid.UUID `json:"id"`
	Alias     string    `json:"alias"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type Original struct {
	Original string `json:"original"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
