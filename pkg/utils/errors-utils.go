package utils

import "encoding/json"

func ResponseError(message string, code int) []byte {
	if code == 0 {
		code = 400
	}

	respError := map[string]interface{}{
		"code":    code,
		"message": message,
	}
	respBody, _ := json.Marshal(respError)
	return respBody
}
