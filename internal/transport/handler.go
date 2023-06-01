package transport

import (
	"context"
	"encoding/json"
	chi "github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router  *chi.Mux
	Service Service
	Server  *http.Server
}

type Response struct {
	Message string `json:"message"`
}

func NewHandler(service Service) *Handler {
	h := &Handler{

		Router:  chi.NewRouter(),
		Service: service,
	}
	h.mapRoutes()
	h.Server = &http.Server{
		Addr:         "0.0.0.0:8081",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.Get("/profiles/{id}", JWTAuth(h.GetProfile))
	h.Router.Get("/profiles/user/{id}", JWTAuth(h.GetProfileByUser))
	h.Router.Get("/profiles", JWTAuth(h.GetProfiles))
	h.Router.Post("/profiles", JWTAuth(h.CreateProfile))
}

func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive!"}); err != nil {
		log.Errorf("Error getting profile: %v", err)
	}
}

// Serve - gracefully serves our newly set up handler function
func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shutting down gracefully")
	return nil
}
