package api

import (
	"cloud.google.com/go/civil"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
		Addr:    ":3000",
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
	s.router.Get("/v1/habits", s.getHabits)
	s.router.Get("/v1/habit/{habitID}", s.getHabit)
	s.router.Post("/v1/habit", s.postHabit)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{OK: true}
	err := render.Render(w, r, health)
	if err != nil {
		return
	}
}

func (s *Server) getHabits(w http.ResponseWriter, r *http.Request) {
	health := HabitsResponse{
		Habits: []Habit{
			{
				Description: "Study Dutch",
				Days: []civil.Date{
					{Day: 1, Month: time.April, Year: 2023},
					{Day: 2, Month: time.April, Year: 2023},
					{Day: 3, Month: time.April, Year: 2023},
				},
			},
		}}
	err := render.Render(w, r, health)
	if err != nil {
		return
	}
}

func (s *Server) postHabit(w http.ResponseWriter, r *http.Request) {
	var habit Habit
	json.NewDecoder(r.Body).Decode(&habit)
	return
}

func (s *Server) getHabit(w http.ResponseWriter, r *http.Request) {
	habitID := chi.URLParam(r, "habitID")
	w.Write([]byte(fmt.Sprintf("something something: %s", habitID)))
	return
}
