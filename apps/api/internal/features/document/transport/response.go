package transport

import "time"

type UploadResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"fileName"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
