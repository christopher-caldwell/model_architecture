package health

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{Message: "ready"})
}
