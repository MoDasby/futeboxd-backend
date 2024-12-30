package cookies

import (
	"net/http"
)

func CreateSessionCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "sid",
		Value:    token,
		Path:     "/",
		MaxAge:   31 * 24 * 60 * 60,
		HttpOnly: true,
	}
}

func DeleteSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
}
