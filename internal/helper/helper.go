package helper

func ResponseFormat(code int, message string, data any) map[string]any {
	var result = make(map[string]any)
	result["code"] = code
	result["message"] = message
	if data != nil {
		result["data"] = data
	}
	return result
}

func ResponseGetOrderFormat(code int, message string, userID int, data any) map[string]any {
	var result = make(map[string]any)
	result["code"] = code
	result["message"] = message
	result["id"] = userID
	if data != nil {
		result["data"] = data
	}
	return result
}
