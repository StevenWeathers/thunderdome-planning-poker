package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/StevenWeathers/thunderdome-planning-poker/api"
	"github.com/spf13/viper"
)

//go:embed dist
var f embed.FS

func getFileSystem(useOS bool) (http.FileSystem, fs.FS) {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("dist")), fs.FS(os.DirFS("dist"))
	}

	fsys, err := fs.Sub(f, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys), fs.FS(fsys)
}

func (s *server) routes() {
	HFS, FSS := getFileSystem(embedUseOS)
	staticHandler := http.FileServer(HFS)

	// api (used by the webapp but can be enabled for external use)
	apiConfig := &api.ApiConfig{
		AppDomain:          s.config.AppDomain,
		FrontendCookieName: s.config.FrontendCookieName,
		SecureCookieName:   viper.GetString("http.backend_cookie_name"),
		SecureCookieFlag:   viper.GetBool("http.secure_cookie"),
		PathPrefix:         s.config.PathPrefix,
		ExternalAPIEnabled: s.config.ExternalAPIEnabled,
	}
	api.Init(apiConfig, s.router, s.database, s.email, s.cookie)

	// static assets
	s.router.PathPrefix("/static/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	s.router.PathPrefix("/img/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	s.router.PathPrefix("/lang/").Handler(http.StripPrefix(s.config.PathPrefix, staticHandler))
	// user avatar generation
	if s.config.AvatarService == "goadorable" || s.config.AvatarService == "govatar" {
		s.router.PathPrefix("/avatar/{width}/{id}/{avatar}").Handler(s.handleUserAvatar()).Methods("GET")
		s.router.PathPrefix("/avatar/{width}/{id}").Handler(s.handleUserAvatar()).Methods("GET")
	}

	// handle index.html
	s.router.PathPrefix("/").HandlerFunc(s.handleIndex(FSS))
}
