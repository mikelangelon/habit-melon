package postgres

import (
	"cloud.google.com/go/civil"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mikelangelon/habit-melon/internal/app"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateHabit(ctx context.Context, habit app.Habit) (app.Habit, error) {
	q := `INSERT INTO habits (description) VALUES ($1) RETURNING habit_id;`
	var id int64
	err := r.db.QueryRowContext(ctx, q, habit.Description).Scan(&id)
	if err != nil {
		return app.Habit{}, err
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

func (r *Repository) DeleteHabitDays(ctx context.Context, id int64, days *[]civil.Date) error {
	if days == nil {
		q := `DELETE FROM habit_days WHERE habit_id = $1`
		_, err := r.db.ExecContext(ctx, q, id)
		if err != nil {
			return fmt.Errorf("problem deleting day: %w", err)
		}
		return nil
	}
	for _, v := range *days {
		q := `DELETE FROM habit_days WHERE habit_id = $1 AND day = $2`
		_, err := r.db.ExecContext(ctx, q, id, v.String())
		if err != nil {
			return fmt.Errorf("problem deleting day: %w", err)
		}
	}
	return nil
}

func (r *Repository) UpdateHabit(ctx context.Context, habit app.Habit) (app.Habit, error) {
	q := `UPDATE habits SET description = $1 WHERE habit_id = $2;`
	_, err := r.db.ExecContext(ctx, q, habit.Description, habit.ID)
	if err != nil {
		return app.Habit{}, err
	}
	err = r.DeleteHabitDays(ctx, *habit.ID, nil)
	if err != nil {
		return app.Habit{}, err
	}
	if err := r.CreateHabitDays(ctx, *habit.ID, habit.Days); err != nil {
		return habit, err
	}
	return habit, nil
}

func (r *Repository) DeleteHabit(ctx context.Context, habitID int64) error {
	q := `DELETE FROM habits WHERE habit_id = $1`
	_, err := r.db.ExecContext(ctx, q, habitID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetHabit(ctx context.Context, habitID int64) (app.Habit, error) {
	var habit app.Habit
	habit.ID = &habitID
	q := `SELECT description FROM habits WHERE habit_id = $1`
	err := r.db.QueryRowContext(ctx, q, habitID).Scan(&habit.Description)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return habit, app.ErrHabitNotFound
		}
		return habit, err
	}

	q = `SELECT day FROM habit_days WHERE habit_id = $1`
	rows, err := r.db.QueryContext(ctx, q, habitID)
	var times []civil.Date
	for rows.Next() {
		var t time.Time
		err := rows.Scan(&t)
		if err != nil {
			return habit, err
		}
		times = append(times, civil.DateOf(t))
	}
	habit.Days = times
	if err != nil {
		return habit, err
	}
	return habit, nil
}

func (r *Repository) GetAllHabit(ctx context.Context) ([]app.Habit, error) {
	var habits []app.Habit
	q := `SELECT habit_id, description FROM habits`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return habits, err
	}
	for rows.Next() {
		var habit app.Habit
		err := rows.Scan(&habit.ID, &habit.Description)
		if err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}
	if err != nil {
		return habits, err
	}
	return habits, nil
}
