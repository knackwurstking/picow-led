package picow

type Response[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
	ID    ID     `json:"id"`
}
