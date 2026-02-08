package handler

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
