package api

import (
	"net/http"
)

type healthResponse struct {
	OK bool `json:"ok"`
}

func (hr healthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
