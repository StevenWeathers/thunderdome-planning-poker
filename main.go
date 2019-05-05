package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "thunderdome/statik"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
)

func main() {
	var listenPort = fmt.Sprintf(":%s", GetEnv("PORT", "8080"))
	go h.run()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	staticHandler := http.FileServer(statikFS)
	
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(staticHandler)
	router.PathPrefix("/js/").Handler(staticHandler)
	router.PathPrefix("/img/").Handler(staticHandler)
	router.HandleFunc("/api/user", RegisterUserHandler).Methods("POST")
	router.HandleFunc("/api/battle", CreateBattleHandler).Methods("POST")
	router.HandleFunc("/api/battle/{id}", serveWs)
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