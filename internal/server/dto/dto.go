package dto

type ErrorResponse[T string | map[string]string] struct {
	Error T `json:"error"`
}
