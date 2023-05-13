package repository

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type HabitPostgres struct {
	dbpool *pgxpool.Pool
}

func NewHabitPostgres(dbpool *pgxpool.Pool) *HabitPostgres {
	return &HabitPostgres{dbpool: dbpool}
}

func (r *HabitPostgres) Create(userId int, habit models.Habit) (int, error) {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	var habitId int
	// create a habit
	createHabitQuery := "INSERT INTO habit (title, description) VALUES ($1, $2) RETURNING id"
	rowHabit := tx.QueryRow(context.Background(), createHabitQuery, habit.Title, habit.Description)
	if err := rowHabit.Scan(&habitId); err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	// create an empty tracker for a habit
	var trackerId int
	createHabitTrackerQuery := "INSERT INTO habit_tracker (habit_id) VALUES ($1) RETURNING id"
	rowTracker := tx.QueryRow(context.Background(), createHabitTrackerQuery, userId)
	err = rowTracker.Scan(&trackerId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	// link habit to a user and a tracker to a habit
	createUsersHabitsQuery := "INSERT INTO user_habit (user_id, habit_id, habit_tracker_id) VALUES ($1, $2, $3)"
	_, err = tx.Exec(context.Background(), createUsersHabitsQuery, userId, habitId, trackerId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	return habitId, tx.Commit(context.Background())
}

func (r *HabitPostgres) GetAll(userId int) ([]models.Habit, error) {
	var habits []models.Habit
	query := "SELECT tl.id, tl.title, tl.description FROM habit tl INNER JOIN user_habit ul on tl.id = ul.habit_id WHERE ul.user_id = $1"
	// err := r.db.QueryRow(context.Background(), query, userId).Scan(&habits)

	rowsHabits, err := r.dbpool.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return habits, err
	}

	defer rowsHabits.Close()

	for rowsHabits.Next() {
		habits, err = pgx.CollectRows(rowsHabits, pgx.RowToStructByName[models.Habit])
		if err != nil {
			fmt.Fprintf(os.Stderr, "rowsHabits CollectRows failed: %v\n", err)
			return habits, err
		}
	}

	// for rowHabit.Scan(
	// 	&habits.Id,
	// 	&habits.Title,
	// 	&habits.Description,
	// )
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	return habit, err
	// }

	return habits, err
}

func (r *HabitPostgres) GetById(userId, habitId int) (models.Habit, error) {
	var habit models.Habit

	query := "SELECT tl.id, tl.title, tl.description FROM habit tl INNER JOIN user_habit ul on tl.id = ul.habit_id WHERE ul.user_id = $1 AND ul.habit_id = $2"

	// err := r.db.QueryRow(context.Background(), query, userId, habitId).Scan(&habit)

	rowHabit := r.dbpool.QueryRow(context.Background(), query, userId, habitId)

	err := rowHabit.Scan(
		&habit.Id,
		&habit.Title,
		&habit.Description,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return habit, err
	}

	return habit, err
}

func (r *HabitPostgres) Delete(userId, habitId int) error {
	queryTracker := "DELETE FROM habit_tracker tl USING user_habit ul WHERE tl.habit_id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2"
	_, err := r.dbpool.Exec(context.Background(), queryTracker, userId, habitId)
	if err != nil {
		return err
	}

	query := "DELETE FROM habit tl USING user_habit ul WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2"
	_, err = r.dbpool.Exec(context.Background(), query, userId, habitId)

	return err
}

func (r *HabitPostgres) Update(userId, habitId int, input models.UpdateHabitInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE habit tl SET %s FROM user_habit ul WHERE tl.id = ul.habit_id AND ul.habit_id=$%d AND ul.user_id=$%d",
		setQuery, argId, argId+1)

	args = append(args, habitId, userId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.dbpool.Exec(context.Background(), query, args...)
	return err
}
