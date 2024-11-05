package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Set the "token" cookie with an expired date to delete it
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		Expires:  time.Unix(0, 0), // Expire immediately
		MaxAge:   -1,              // Alternatively, set MaxAge to -1 for immediate deletion
	})

	resp := Response{}
	resp.Code = make(map[int]string)
	resp.Code[http.StatusOK] = "OK"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
