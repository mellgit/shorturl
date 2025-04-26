package redirect

import "time"

type Click struct {
	ID        int64     `json:"id"`
	Alias     string    `json:"alias"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
