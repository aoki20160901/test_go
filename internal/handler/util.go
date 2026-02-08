package handler

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // 余計なフィールド防止（おすすめ）
	return decoder.Decode(v)
}

func RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
