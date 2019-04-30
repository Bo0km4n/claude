package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ping := struct {
			Message string `json:"message"`
		}{
			Message: "ok",
		}

		res, _ := json.Marshal(ping)

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.ListenAndServe(":8080", nil)
}
