package repo

type Result[T any] struct {
	Code    int64  `json:"code"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
