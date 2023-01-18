package transport

type HTTPBaseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
