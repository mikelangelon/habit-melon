package api

import (
	"cloud.google.com/go/civil"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/mikelangelon/habit-melon/internal/app"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	router *chi.Mux

	habitStore   map[string]Habit
	habitService habitService
}

type habitService interface {
	GetAll(ctx context.Context) ([]app.Habit, error)
	Get(ctx context.Context, habitID int64) (app.Habit, error)
	Create(ctx context.Context, habit app.Habit) (app.Habit, error)
	Update(ctx context.Context, habit app.Habit) (app.Habit, error)
	Delete(ctx context.Context, habitID int64) error
}

func NewServer(habitService habitService) *Server {
	srv := &Server{
		router:       chi.NewRouter(),
		habitService: habitService,
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
	s.router.Put("/v1/habit/{habitID}", s.putHabit)
	s.router.Delete("/v1/habit/{habitID}", s.deleteHabit)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{OK: true}
	err := render.Render(w, r, health)
	if err != nil {
		return
	}
}

func (s *Server) getHabits(w http.ResponseWriter, r *http.Request) {
	all, err := s.habitService.GetAll(context.TODO())
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := habitsOnResponse(all, w); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	return
}

func (s *Server) getHabit(w http.ResponseWriter, r *http.Request) {
	habitID := chi.URLParam(r, "habitID")
	id, err := strconv.ParseInt(habitID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}
	habit, err := s.habitService.Get(context.TODO(), id)
	if err != nil {
		switch err {
		case app.ErrHabitNotFound:
			http.Error(w, http.StatusText(404), 404)
		default:
			http.Error(w, http.StatusText(500), 500)
		}

		return
	}
	if err := habitOnResponse(habit, w); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	return
}
func (s *Server) postHabit(w http.ResponseWriter, r *http.Request) {
	var habit app.Habit
	err := json.NewDecoder(r.Body).Decode(&habit)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	habit, err = s.habitService.Create(context.TODO(), habit)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := habitOnResponse(habit, w); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
func (s *Server) putHabit(w http.ResponseWriter, r *http.Request) {
	habitID := chi.URLParam(r, "habitID")
	id, err := strconv.ParseInt(habitID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}
	var habit app.Habit
	err = json.NewDecoder(r.Body).Decode(&habit)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	habit.ID = &id
	habit, err = s.habitService.Update(context.TODO(), habit)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := habitOnResponse(habit, w); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	return
}
func (s *Server) deleteHabit(w http.ResponseWriter, r *http.Request) {
	habitID := chi.URLParam(r, "habitID")
	id, err := strconv.ParseInt(habitID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}
	err = s.habitService.Delete(context.TODO(), id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	return
}

func habitOnResponse(habit app.Habit, w http.ResponseWriter) error {
	bytes, err := json.Marshal(habit)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func habitsOnResponse(habits []app.Habit, w http.ResponseWriter) error {
	bytes, err := json.Marshal(habits)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
