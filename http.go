package main

import (
	"bytes"
	"embed"
	"html/template"
	"image"
	"image/png"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	api "github.com/StevenWeathers/thunderdome-planning-poker/api"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gorilla/mux"
	"github.com/ipsn/go-adorable"
	"github.com/o1egl/govatar"
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
	apiConfig := &api.Config{
		AppDomain:          s.config.AppDomain,
		FrontendCookieName: s.config.FrontendCookieName,
		SecureCookieName:   viper.GetString("http.backend_cookie_name"),
		SecureCookieFlag:   viper.GetBool("http.secure_cookie"),
		SessionCookieName:  viper.GetString("http.session_cookie_name"),
		PathPrefix:         s.config.PathPrefix,
		ExternalAPIEnabled: s.config.ExternalAPIEnabled,
		UserAPIKeyLimit:    s.config.UserAPIKeyLimit,
		LdapEnabled:        s.config.LdapEnabled,
		FeaturePoker:       viper.GetBool("feature.poker"),
		FeatureRetro:       viper.GetBool("feature.retro"),
		FeatureStoryboard:  viper.GetBool("feature.storyboard"),
	}
	api.Init(apiConfig, s.router, s.db, s.email, s.cookie)

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

// get the index template from embedded filesystem
func (s *server) getIndexTemplate(FSS fs.FS) *template.Template {
	// get the html template from dist, have it ready for requests
	tmplContent, ioErr := fs.ReadFile(FSS, "static/index.html")
	if ioErr != nil {
		log.Println("Error opening index template")
		if !embedUseOS {
			log.Fatal(ioErr)
		}
	}

	tmplString := string(tmplContent)
	tmpl, tmplErr := template.New("index").Parse(tmplString)
	if tmplErr != nil {
		log.Println("Error parsing index template")
		if !embedUseOS {
			log.Fatal(tmplErr)
		}
	}

	return tmpl
}

/*
	Handlers
*/

// handleIndex parses the index html file, injecting any relevant data
func (s *server) handleIndex(FSS fs.FS) http.HandlerFunc {
	type AppConfig struct {
		AllowedPointValues    []string
		DefaultPointValues    []string
		ShowWarriorRank       bool
		AvatarService         string
		ToastTimeout          int
		AllowGuests           bool
		AllowRegistration     bool
		AllowJiraImport       bool
		DefaultLocale         string
		FriendlyUIVerbs       bool
		AppVersion            string
		CookieName            string
		PathPrefix            string
		ExternalAPIEnabled    bool
		UserAPIKeyLimit       int
		CleanupGuestsDaysOld  int
		CleanupBattlesDaysOld int
		CleanupRetrosDaysOld  int
		ShowActiveCountries   bool
		LdapEnabled           bool
		FeaturePoker          bool
		FeatureRetro          bool
		FeatureStoryboard     bool
	}
	type UIConfig struct {
		AnalyticsEnabled bool
		AnalyticsID      string
		AppConfig        AppConfig
		ActiveAlerts     []interface{}
	}

	tmpl := s.getIndexTemplate(FSS)

	appConfig := AppConfig{
		AllowedPointValues:    viper.GetStringSlice("config.allowedPointValues"),
		DefaultPointValues:    viper.GetStringSlice("config.defaultPointValues"),
		ShowWarriorRank:       viper.GetBool("config.show_warrior_rank"),
		AvatarService:         viper.GetString("config.avatar_service"),
		ToastTimeout:          viper.GetInt("config.toast_timeout"),
		AllowGuests:           viper.GetBool("config.allow_guests"),
		AllowRegistration:     viper.GetBool("config.allow_registration") && viper.GetString("auth.method") == "normal",
		AllowJiraImport:       viper.GetBool("config.allow_jira_import"),
		DefaultLocale:         viper.GetString("config.default_locale"),
		FriendlyUIVerbs:       viper.GetBool("config.friendly_ui_verbs"),
		ExternalAPIEnabled:    s.config.ExternalAPIEnabled,
		UserAPIKeyLimit:       s.config.UserAPIKeyLimit,
		AppVersion:            s.config.Version,
		CookieName:            s.config.FrontendCookieName,
		PathPrefix:            s.config.PathPrefix,
		CleanupGuestsDaysOld:  viper.GetInt("config.cleanup_guests_days_old"),
		CleanupBattlesDaysOld: viper.GetInt("config.cleanup_battles_days_old"),
		CleanupRetrosDaysOld:  viper.GetInt("config.cleanup_retros_days_old"),
		ShowActiveCountries:   viper.GetBool("config.show_active_countries"),
		LdapEnabled:           s.config.LdapEnabled,
		FeaturePoker:          viper.GetBool("feature.poker"),
		FeatureRetro:          viper.GetBool("feature.retro"),
		FeatureStoryboard:     viper.GetBool("feature.storyboard"),
	}

	data := UIConfig{
		AnalyticsEnabled: s.config.AnalyticsEnabled,
		AnalyticsID:      s.config.AnalyticsID,
		AppConfig:        appConfig,
	}

	api.ActiveAlerts = s.db.GetActiveAlerts() // prime the active alerts cache

	return func(w http.ResponseWriter, r *http.Request) {
		data.ActiveAlerts = api.ActiveAlerts // get latest alerts from memory

		if embedUseOS {
			tmpl = s.getIndexTemplate(FSS)
		}

		tmpl.Execute(w, data)
	}
}

// handleUserAvatar creates an avatar for the given user by ID
func (s *server) handleUserAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		Width, _ := strconv.Atoi(vars["width"])
		UserID := vars["id"]
		AvatarGender := govatar.MALE
		userGender, ok := vars["avatar"]
		if ok {
			if userGender == "female" {
				AvatarGender = govatar.FEMALE
			}
		}

		var avatar image.Image
		if s.config.AvatarService == "govatar" {
			avatar, _ = govatar.GenerateForUsername(AvatarGender, UserID)
		} else { // must be goadorable
			var err error
			avatar, _, err = image.Decode(bytes.NewReader(adorable.PseudoRandom([]byte(UserID))))
			if err != nil {
				log.Fatalln(err)
			}
		}

		img := transform.Resize(avatar, Width, Width, transform.Linear)
		buffer := new(bytes.Buffer)

		if err := png.Encode(buffer, img); err != nil {
			log.Println("unable to encode image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
