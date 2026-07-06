package http

type Envelope[T any] struct {
	Data T `json:"data"`
}
