package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"time"

	_ "thunderdome/statik"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

// getEnv gets environment variable matching key string
// and if it finds none uses fallback string
// returning either the matching or fallback string
func getEnv(key string, fallback string) string {
	var result = os.Getenv(key)

	if result == "" {
		result = fallback
	}

	return result
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	var listenPort = fmt.Sprintf(":%s", getEnv("PORT", "8080"))

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	staticHandler := http.FileServer(statikFS)
	
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(staticHandler)
	router.PathPrefix("/js/").Handler(staticHandler)
	router.HandleFunc("/echo", echo)
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