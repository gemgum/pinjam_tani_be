package responses

type MapResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONWebResponse(code int, status string, msg string, data interface{}) MapResponse {
	return MapResponse{
		Code:    code,
		Status:  status,
		Message: msg,
		Data:    data,
	}
}
