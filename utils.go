package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

// GetEnv gets environment variable matching key string
// and if it finds none uses fallback string
// returning either the matching or fallback string
func GetEnv(key string, fallback string) string {
	var result = os.Getenv(key)

	if result == "" {
		result = fallback
	}

	return result
}

func GetIntEnv(key string, fallback int) int {
	var intResult = fallback
	var stringResult = os.Getenv(key)

	if stringResult != "" {
		v, _ := strconv.Atoi(stringResult)
		intResult = v
	}

	return intResult
}

// RespondWithJSON takes a payload and writes the response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
