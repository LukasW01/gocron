package http

import (
	"encoding/json"
	"gocron/src"
	"log"
	"net/http"
	"os"
)

func Server(port string) {
	log.Println("Opening port", port, "for health checking")

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(cron.Proc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cron.Proc.Status.ExitStatus != 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if _, err := w.Write(json); err != nil {
		return
	}
}
