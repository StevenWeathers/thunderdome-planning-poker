package cookie

import "github.com/gorilla/securecookie"

type Cookie struct {
	config Config
	sc     *securecookie.SecureCookie
}

type Config struct {
	// the domain of the application for cookie securing
	AppDomain string
	// PathPrefix allows the application to be run on a shared domain
	PathPrefix string
	// CookieHashKey
	CookieHashKey string
	// name of the cookie used exclusively by the UI
	FrontendCookieName string
	// name of the user cookie
	SecureCookieName string
	// name of the user session cookie used for authenticated sessions
	SessionCookieName string
	// controls whether the cookie is set to secure, only works over HTTPS
	SecureCookieFlag bool
}
