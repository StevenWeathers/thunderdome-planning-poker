package cookie

import (
	"fmt"
	"net/http"
	"time"
)

// CreateAuthStateCookie creates the oauth2 state validation cookie
func (s *Cookie) CreateAuthStateCookie(w http.ResponseWriter, state string) error {
	return s.CreateCookie(w, s.config.AuthStateCookieName, state, int(time.Minute.Seconds()*10))
}

// ValidateAuthStateCookie retrieves the authstate cookie and validates its value against the state from
// auth provider and deletes the authstate cookie
func (s *Cookie) ValidateAuthStateCookie(w http.ResponseWriter, r *http.Request, state string) error {
	cookieVal, err := s.GetCookie(w, r, s.config.AuthStateCookieName)
	if err != nil {
		return err
	}

	_ = s.DeleteAuthStateCookie(w)

	if cookieVal != state {
		return fmt.Errorf("INVALID_AUTH_STATE")
	}

	return nil
}

// DeleteAuthStateCookie deletes the oauth2 state validation cookie
func (s *Cookie) DeleteAuthStateCookie(w http.ResponseWriter) error {
	return s.CreateCookie(w, s.config.AuthStateCookieName, "", -1)
}
