package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"workout-microservice/internal/validator"
)

const insertPrQuery = `INSERT INTO exercise_prs(USER_ID, EXERCISE_ID, PR) VALUES ($1, $2, $3)`
const selectPrQuery = `SELECT pr FROM exercise_prs WHERE (user_id, exercise_id) = ($1, $2)`
const updatePrQuery = `UPDATE exercise_prs SET pr = $1 WHERE (user_id, exercise_id) = ($1, $2)`
const deletePrQuery = `DELETE FROM exercise_prs WHERE (user_id, exercise_id) = ($1, $2)`

type Pr struct {
	UserId         int `json:"user_id"`
	ExerciseId     int `json:"exercise_id"`
	PersonalRecord int `json:"personal_record"`
}

type PrModel struct {
	db *sql.DB
}

func (p PrModel) Insert(pr Pr) error {
	oldPr, err := p.Get(pr.UserId, pr.ExerciseId)
	if err != nil {
		// if there aren't any rows then we insert the row
		if errors.Is(err, sql.ErrNoRows) {
			err2 := p.runQuery(pr, insertExerciseQuery)
			if err2 != nil {
				fmt.Printf("error while inserting row with user id: %d and exercise id: %d \n",
					pr.UserId,
					pr.ExerciseId)
				return err2
			}

			return nil
		}
	}

	err = p.Update(*oldPr)
	if err != nil {
		return err
	}

	return nil
}

func (p PrModel) Update(pr Pr) error {
	return p.runQuery(pr, updatePrQuery)
}

func (p PrModel) Delete(pr Pr) error {
	return p.runQuery(pr, deletePrQuery)
}

func (p PrModel) runQuery(pr Pr, query string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	args := []interface{}{pr.UserId, pr.ExerciseId}
	_, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p PrModel) Get(userId int, exerciseId int) (*Pr, error) {
	pr := Pr{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	args := []interface{}{userId, exerciseId}
	err := p.db.QueryRowContext(ctx, selectPrQuery, args...).Scan(&pr.PersonalRecord)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Printf("no pr for user id: %d and exercise id: %d \n", userId, exerciseId)
			return nil, err
		default:
			return nil, err
		}
	}

	pr.UserId = userId
	pr.ExerciseId = exerciseId

	return &pr, nil
}

func ValidatePr(v *validator.Validator, pr *Pr) {
	v.Check(pr.UserId >= 1, "User id", " must be >= 1")
	v.Check(pr.ExerciseId >= 1, "Exercise id", " must be >= 1")
}
