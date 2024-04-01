package cookie

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

var re = regexp.MustCompile(`%(2[346BF]|3[AC-F]|40|5[BDE]|60|7[BCD])`)

func unescapeRegexpStringMatch(s string) string {
	a, _ := url.PathUnescape(s)
	return a
}

func createJsonCookieValue(value any) (string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	// js-cookie reads json cookie values using value.replace(/(%[\dA-F]{2})+/gi, decodeURIComponent)
	// js-cookie writes json cookie values using encodeURIComponent(value).replace(
	//      /%(2[346BF]|3[AC-F]|40|5[BDE]|60|7[BCD])/g,
	//      decodeURIComponent
	//    )
	encodedValue := url.PathEscape(string(jsonValue))
	s := re.ReplaceAllStringFunc(encodedValue, unescapeRegexpStringMatch)

	return s, nil
}

// CreateUserCookie creates the users Cookie
func (s *Cookie) CreateUserCookie(w http.ResponseWriter, UserID string) error {
	return s.CreateCookie(w, s.config.SecureCookieName, UserID, int(time.Hour.Seconds()*24*365))
}

// CreateSessionCookie creates the user's session Cookie
func (s *Cookie) CreateSessionCookie(w http.ResponseWriter, SessionID string) error {
	return s.CreateCookie(w, s.config.SessionCookieName, SessionID, int(time.Hour.Seconds()*24*30))
}

// CreateUserUICookie creates the user's frontend UI cookie
func (s *Cookie) CreateUserUICookie(w http.ResponseWriter, userUiCookie thunderdome.UserUICookie) error {
	encodedValue, err := createJsonCookieValue(userUiCookie)
	if err != nil {
		return fmt.Errorf("error creating encoded json for cookie: %w", err)
	}

	c := &http.Cookie{
		Name:     s.config.FrontendCookieName,
		Value:    encodedValue,
		Path:     s.config.PathPrefix + "/",
		Domain:   s.config.AppDomain,
		MaxAge:   int(time.Hour.Seconds() * 24 * 365),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, c)

	return nil
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
