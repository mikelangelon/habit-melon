package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	srv := &Server{
		router: chi.NewRouter(),
	}

	srv.routes()
	return srv
}

func (s *Server) Start() {
	server := http.Server{
		Addr:    ":8100",
		Handler: s.router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("http.ListenAndServe failed: %v\n", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.Use(render.SetContentType(render.ContentTypeJSON))
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	s.router.Get("/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{OK: true}
	err := render.Render(w, r, health)
	if err != nil {
		return
	}
}
