package app

import (
	"cloud.google.com/go/civil"
	"context"
	"errors"
)

var ErrHabitNotFound = errors.New("habit not found")

type Habit struct {
	ID          *int64       `json:"id"`
	Description string       `json:"description"`
	Days        []civil.Date `json:"days"`
}
type HabitService struct {
	repo repo
}

type repo interface {
	CreateHabit(ctx context.Context, habit Habit) (Habit, error)
	UpdateHabit(ctx context.Context, habit Habit) (Habit, error)
	GetHabit(ctx context.Context, habitID int64) (Habit, error)
	GetAllHabit(ctx context.Context) ([]Habit, error)
	DeleteHabit(ctx context.Context, habitID int64) error
}

func NewHabitService(repo repo) HabitService {
	return HabitService{
		repo: repo,
	}
}

func (h HabitService) Create(ctx context.Context, habit Habit) (Habit, error) {
	return h.repo.CreateHabit(ctx, habit)
}

func (h HabitService) Update(ctx context.Context, habit Habit) (Habit, error) {
	return h.repo.UpdateHabit(ctx, habit)
}

func (h HabitService) Delete(ctx context.Context, habitID int64) error {
	return h.repo.DeleteHabit(ctx, habitID)
}

func (h HabitService) GetAll(ctx context.Context) ([]Habit, error) {
	return h.repo.GetAllHabit(ctx)
}

func (h HabitService) Get(ctx context.Context, habitID int64) (Habit, error) {
	return h.repo.GetHabit(ctx, habitID)
}
