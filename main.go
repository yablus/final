package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/yablus/final/internal/handlers"
	"github.com/yablus/final/internal/service"
)

// http://127.0.0.1:8282/api

func main() {
	c := handleConnection()
	http.ListenAndServe(":8282", c)
}

func handleConnection() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})
	h := &handlers.ServiceHandler{Data: service.NewService()}
	r.Mount("/api", MountRoute(h))
	r.Mount("/test", MountRouteTests(h))
	return r
}

func MountRoute(h *handlers.ServiceHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetData) // GET /api
	//r.Get("/", test.HandleTest) // GET /api *для тестирования веб-сервиса
	return r
}

func MountRouteTests(h *handlers.ServiceHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetData)                     // GET /test
	r.Get("/ResultSetT", h.GetResultSetTData) // GET /test/ResultSetT
	r.Get("/SMS", h.GetSMSData)               // GET /test/SMS
	return r
}
