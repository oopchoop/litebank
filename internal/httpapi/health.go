package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func HealthHandler(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

		defer cancel()

		w.Header().Set("Content-type", "application/json")

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)

			json.NewEncoder(w).Encode(map[string]string{"status": "error"})

			log.Printf("Ошибка %v:", err)
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
