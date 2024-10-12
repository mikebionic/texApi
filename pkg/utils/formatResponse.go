package utils

type UniversalResponse struct {
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:"errorMsg"`
}

func FormatResponse(message string, data interface{}) UniversalResponse {
	return UniversalResponse{
		Message:  message,
		Success:  true,
		Data:     data,
		ErrorMsg: "",
	}
}

func FormatErrorResponse(message string, errorMsg string) UniversalResponse {
	if errorMsg == "" {
		errorMsg = message
	}
	return UniversalResponse{
		Message:  message,
		Success:  false,
		Data:     nil,
		ErrorMsg: errorMsg,
	}
}
