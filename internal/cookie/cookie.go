package cookie

import "github.com/gorilla/securecookie"

func New(config Config) *Cookie {
	c := Cookie{
		config: config,
	}

	c.sc = securecookie.New([]byte(c.config.CookieHashKey), nil)

	return &c
}
