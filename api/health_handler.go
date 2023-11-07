package api

import (
	"net/http"
)

type HealthResponse struct {
	OK bool `json:"ok"`
}

func (hr HealthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
