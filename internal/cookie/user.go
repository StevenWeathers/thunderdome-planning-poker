package cookie

import (
	"fmt"
	"net/http"
)

// CreateUserCookie creates the users Cookie
func (s *Cookie) CreateUserCookie(w http.ResponseWriter, UserID string) error {
	encoded, err := s.sc.Encode(s.config.SecureCookieName, UserID)
	if err != nil {
		return err

	}

	cookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    encoded,
		Path:     s.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   86400 * 365,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return nil
}

// CreateSessionCookie creates the user's session Cookie
func (s *Cookie) CreateSessionCookie(w http.ResponseWriter, SessionID string) error {
	encoded, err := s.sc.Encode(s.config.SessionCookieName, SessionID)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     s.config.SessionCookieName,
		Value:    encoded,
		Path:     s.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   86400 * 30,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	return nil
}

// ClearUserCookies wipes the frontend and backend cookies
// used in the event of bad Cookie reads
func (s *Cookie) ClearUserCookies(w http.ResponseWriter) {
	feCookie := &http.Cookie{
		Name:   s.config.FrontendCookieName,
		Value:  "",
		Path:   s.config.PathPrefix + "/",
		MaxAge: -1,
	}
	beCookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    "",
		Path:     s.config.PathPrefix + "/",
		Domain:   s.config.AppDomain,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}
	seCookie := &http.Cookie{
		Name:     s.config.SessionCookieName,
		Value:    "",
		Path:     s.config.PathPrefix + "/",
		Domain:   s.config.AppDomain,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, feCookie)
	http.SetCookie(w, beCookie)
	http.SetCookie(w, seCookie)
}

// ValidateUserCookie returns the UserID from secure cookies or errors if failures getting it
func (s *Cookie) ValidateUserCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var UserID string

	if cookie, err := r.Cookie(s.config.SecureCookieName); err == nil {
		var value string
		if err = s.sc.Decode(s.config.SecureCookieName, cookie.Value, &value); err == nil {
			UserID = value
		} else {
			s.ClearUserCookies(w)
			return "", fmt.Errorf("INVALID_USER_COOKIE")
		}
	} else {
		return "", fmt.Errorf("NO_USER_COOKIE")
	}

	return UserID, nil
}

// ValidateSessionCookie returns the SessionID from secure cookies or errors if failures getting it
func (s *Cookie) ValidateSessionCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var SessionID string

	if cookie, err := r.Cookie(s.config.SessionCookieName); err == nil {
		var value string
		if err = s.sc.Decode(s.config.SessionCookieName, cookie.Value, &value); err == nil {
			SessionID = value
		} else {
			s.ClearUserCookies(w)
			return "", fmt.Errorf("INVALID_SESSION_COOKIE")
		}
	} else {
		return "", fmt.Errorf("NO_SESSION_COOKIE")
	}

	return SessionID, nil
}
