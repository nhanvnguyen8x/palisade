package rest

type Response[T any] struct {
	Data T `json:"data"`
}
