package cookie

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

func New(config Config) *Cookie {
	c := Cookie{
		config: config,
	}

	c.sc = securecookie.New([]byte(c.config.CookieHashKey), nil)

	return &c
}

// CreateCookie creates a secure cookie with given cookieName
func (s *Cookie) CreateCookie(w http.ResponseWriter, cookieName string, value string, maxAge int) error {
	encoded, err := s.sc.Encode(cookieName, value)
	if err != nil {
		return err

	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     s.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   maxAge,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return nil
}

// GetCookie returns the value from the cookie or errors if failure getting it
func (s *Cookie) GetCookie(w http.ResponseWriter, r *http.Request, cookieName string) (string, error) {
	var value string

	if cookie, err := r.Cookie(cookieName); err == nil {
		if err = s.sc.Decode(cookieName, cookie.Value, &value); err != nil {
			s.DeleteCookie(w, cookieName)
			return "", fmt.Errorf("INVALID_COOKIE")
		}
	} else {
		return "", fmt.Errorf("COOKIE_NOT_FOUND")
	}

	return value, nil
}

// DeleteCookie deletes the cookie
func (s *Cookie) DeleteCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     s.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   -1,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
