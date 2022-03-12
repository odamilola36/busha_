package dto

import (
	"encoding/json"
	"net/http"
)

// Response is used to shape return json
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

//BuildResponse method to inject data value send out
func BuildResponse(r http.ResponseWriter, response *Response) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(response.Status)

	_ = json.NewEncoder(r).Encode(response)
}
