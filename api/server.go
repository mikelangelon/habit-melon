package api

import (
	"cloud.google.com/go/civil"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"time"
)

type Server struct {
	router *chi.Mux

	habitStore map[string]Habit
}

func NewServer() *Server {
	srv := &Server{
		router: chi.NewRouter(),
		habitStore: map[string]Habit{
			"Study Dutch": {
				Description: "Study Dutch",
				Days: []civil.Date{
					{Day: 1, Month: time.April, Year: 2023},
					{Day: 2, Month: time.April, Year: 2023},
					{Day: 3, Month: time.April, Year: 2023},
				},
			},
		},
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
	habits := HabitsResponse{Habits: convertToList(s.habitStore)}
	err := render.Render(w, r, habits)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func (s *Server) postHabit(w http.ResponseWriter, r *http.Request) {
	var habit Habit
	err := json.NewDecoder(r.Body).Decode(&habit)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	s.habitStore[habit.Description] = habit
	w.WriteHeader(http.StatusCreated)
	return
}

func (s *Server) getHabit(w http.ResponseWriter, r *http.Request) {
	habitID := chi.URLParam(r, "habitID")
	v, ok := s.habitStore[habitID]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	return
}

func convertToList(m map[string]Habit) []Habit {
	var list []Habit
	for _, v := range m {
		list = append(list, v)
	}
	return list
}
