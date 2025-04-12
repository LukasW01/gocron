package cron

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(proc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if proc.Status.ExitStatus != 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if _, err := w.Write(json); err != nil {
		log.Printf("srror writing response: %v", err)
		return
	}
}

func HttpServer(port string) {
	log.Println("Opening port", port, "for health checking")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
