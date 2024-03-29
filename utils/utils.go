package utils

import (
	"encoding/json"
	"net/http"
)

// Message explanation
func Message(status bool, message string) (map[string]interface{}){
	return map[string]interface{} {"status": status, "message": message}
}

// Respond explanation
func Respond(w http.ResponseWriter, data map[string]interface {}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}