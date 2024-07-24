package response

type Response struct {
	Status  uint16                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}
