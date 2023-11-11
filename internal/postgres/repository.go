package postgres

import (
	"cloud.google.com/go/civil"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mikelangelon/habit-melon/internal/core"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateHabit(ctx context.Context, habit core.Habit) (core.Habit, error) {
	q := `INSERT INTO habits (description) VALUES ($1) RETURNING habit_id;`
	var id int64
	err := r.db.QueryRowContext(ctx, q, habit.Description).Scan(&id)
	if err != nil {
		return core.Habit{}, err
	}
	habit.ID = &id
	if err := r.CreateHabitDays(ctx, id, habit.Days); err != nil {
		return habit, err
	}
	return habit, nil
}

func (r *Repository) CreateHabitDays(ctx context.Context, id int64, days []civil.Date) error {
	for _, v := range days {
		q := `INSERT INTO habit_days (habit_id, day) VALUES ($1, $2);`
		_, err := r.db.ExecContext(ctx, q, id, v.String())
		if err != nil {
			return fmt.Errorf("problem inserting day: %w", err)
		}
	}
	return nil
}
