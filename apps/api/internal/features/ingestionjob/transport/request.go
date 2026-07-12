package transport

type UpdateStatusRequest struct {
	Status       string `json:"status" binding:"required"`
	ErrorMessage string `json:"errorMessage"`
}
