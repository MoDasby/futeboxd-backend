package cookies

import (
	"net/http"
	"strings"
	"time"

	"github.com/modasby/futeboxd-backend/core/config"
)

func stringToSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "default":
		return http.SameSiteDefaultMode
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}

func CreateSessionCookie(cfg *config.Cookies, token string, expTime time.Time) *http.Cookie {
	exp := expTime.Sub(time.Now().UTC())

	return &http.Cookie{
		Value:    token,
		Name:     cfg.Name,
		HttpOnly: cfg.HttpOnly,
		Secure:   cfg.Secure,
		SameSite: stringToSameSite(cfg.SameSite),
		Path:     cfg.Path,
		MaxAge:   int(exp.Seconds()),
		Domain:   cfg.Domain,
	}
}

func DeleteSessionCookie(cfg *config.Cookies) *http.Cookie {
	return &http.Cookie{
		Name:   cfg.Name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		Domain: cfg.Domain,
	}
}
