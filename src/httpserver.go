package cron

import (
	"encoding/json"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.MarshalIndent(Current_state, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if Current_state.Last.Exit_status != 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if _, err := w.Write(js); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func Http_server(port string) {
	log.Println("Opening port", port, "for health checking")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
