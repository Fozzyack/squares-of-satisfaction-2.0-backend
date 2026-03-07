package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/env"
)

func SendJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func ErrorJSON(w http.ResponseWriter, errMsg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
}

func DecodeJSON(r *http.Request, out interface{}) error {
	err := json.NewDecoder(r.Body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func CreateCookie(value string, expires_at time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "habbit-tracker",
		Value:    value,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_at,
	}
	if env.GetProduction() {
		cookie.Secure = true
	}
	return cookie

}
