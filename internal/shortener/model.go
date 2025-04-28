package shortener

import "time"

type URL struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Original  string    `json:"original"`
	Alias     string    `json:"alias"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ShortenRequest struct {
	URL      string `json:"url" validate:"required,url"`
	Custom   string `json:"custom"`                            // optional alias
	TTLHours int    `json:"ttl_hours" validate:"gte=1,lte=72"` // optional ttl, default to 24
}

type ErrorResponse struct {
	Error string `json:"error"`
}
