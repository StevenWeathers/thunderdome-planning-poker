package cookie

import (
	"net/http"
	"time"
)

// CreateUserCookie creates the users Cookie
func (s *Cookie) CreateUserCookie(w http.ResponseWriter, UserID string) error {
	return s.CreateCookie(w, s.config.SecureCookieName, UserID, int(time.Hour.Seconds()*24*365))
}

// CreateSessionCookie creates the user's session Cookie
func (s *Cookie) CreateSessionCookie(w http.ResponseWriter, SessionID string) error {
	return s.CreateCookie(w, s.config.SessionCookieName, SessionID, int(time.Hour.Seconds()*24*30))
}

// ClearUserCookies wipes the frontend and backend cookies
// used in the event of bad Cookie reads
func (s *Cookie) ClearUserCookies(w http.ResponseWriter) {
	s.DeleteCookie(w, s.config.SecureCookieName)
	s.DeleteCookie(w, s.config.SessionCookieName)

	feCookie := &http.Cookie{
		Name:   s.config.FrontendCookieName,
		Value:  "",
		Path:   s.config.PathPrefix + "/",
		MaxAge: -1,
	}
	http.SetCookie(w, feCookie)
}

// ValidateUserCookie returns the UserID from secure cookies or errors if failures getting it
func (s *Cookie) ValidateUserCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	return s.GetCookie(w, r, s.config.SecureCookieName)
}

// ValidateSessionCookie returns the SessionID from secure cookies or errors if failures getting it
func (s *Cookie) ValidateSessionCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	return s.GetCookie(w, r, s.config.SessionCookieName)
}
