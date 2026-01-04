package main

import (
	"net/http"

	"github.com/coalaura/plain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var log = plain.New(plain.WithDate(plain.RFC3339Local))

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(log.Middleware())

	r.Post("/report", func(w http.ResponseWriter, r *http.Request) {
		reports := ParseReport(r.Body)
		if len(reports) == 0 {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		for _, report := range reports {
			log.Warnf("CSP [%s] %s blocked %s\n", report.ViolatedDirective, report.DocumentURL, report.BlockedURL)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Listening at http://localhost:9393/")
	log.MustFail(http.ListenAndServe(":9393", r))
}
