package transport

type CreateRequest struct {
	Name string `json:"name" binding:"required"`
}
