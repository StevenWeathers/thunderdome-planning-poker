package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var AppDomain string
var SecureCookieHashkey []byte
var SecureCookieName string = "warriorId"
var Sc = securecookie.New([]byte("some-secret"), nil)

func main() {
	SetupDB() // Sets up DB Connection, and if necessary Tables

	var listenPort = fmt.Sprintf(":%s", GetEnv("PORT", "8080"))
	AppDomain = GetEnv("APP_DOMAIN", "thunderdome.dev")
	SecureCookieHashkey = []byte(GetEnv("COOKIE_HASHKEY", "strongest-avenger"))
	Sc = securecookie.New(SecureCookieHashkey, nil)

	go h.run()

	// box := packr.New("webui", "./dist")
	box := packr.NewBox("./dist")
	staticHandler := http.FileServer(box)

	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(staticHandler)
	router.PathPrefix("/js/").Handler(staticHandler)
	router.PathPrefix("/img/").Handler(staticHandler)
	router.HandleFunc("/api/warrior", RecruitWarriorHandler).Methods("POST")
	router.HandleFunc("/api/battle", CreateBattleHandler).Methods("POST")
	router.HandleFunc("/api/battle/{id}", GetBattleHandler)
	router.HandleFunc("/api/arena/{id}", serveWs)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		staticHandler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    listenPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Access the WebUI via 127.0.0.1" + listenPort)

	log.Fatal(srv.ListenAndServe())
}
