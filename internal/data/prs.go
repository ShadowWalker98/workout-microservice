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

const selectPrQueryByBoth = `SELECT user_id, exercise_prs.exercise_id, exercise_name, exercise_description, pr FROM exercise_prs JOIN exercises
ON exercise_prs.exercise_id = exercises.exercise_id 
WHERE (user_id, exercise_prs.exercise_id) = ($1, $2);`

const selectPrByUserId = `SELECT user_id, exercise_prs.exercise_id, exercise_name, exercise_description, pr FROM exercise_prs JOIN exercises
ON exercise_prs.exercise_id = exercises.exercise_id WHERE user_id = $1;`

const updatePrQuery = `UPDATE exercise_prs SET pr = $1 WHERE (user_id, exercise_id) = ($2, $3)`

const deletePrQuery = `DELETE FROM exercise_prs WHERE (user_id, exercise_id) = ($1, $2)`

type Pr struct {
	UserId         int `json:"user_id"`
	ExerciseId     int `json:"exercise_id"`
	PersonalRecord int `json:"personal_record"`
}

type ConsolidatedPr struct {
	UserId              int    `json:"user_id"`
	ExerciseId          int    `json:"exercise_id"`
	ExerciseName        string `json:"exercise_name"`
	ExerciseDescription string `json:"exercise_description"`
	PersonalRecord      int    `json:"personal_record"`
}

type PrModel struct {
	db *sql.DB
}

func (p PrModel) Insert(pr Pr) error {
	err := checkPr(p.db, pr.UserId, pr.ExerciseId)
	if err != nil {
		// if there aren't any rows then we insert the row
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, ErrRecordNotFound) {
			err2 := p.runQuery(insertPrQuery, []interface{}{pr.UserId, pr.ExerciseId, pr.PersonalRecord})
			if err2 != nil {
				fmt.Printf("error while inserting row with user id: %d and exercise id: %d \n",
					pr.UserId,
					pr.ExerciseId)
				return err2
			}

			return nil
		}
		return err
	}

	err = p.Update(pr)
	if err != nil {
		return err
	}

	return nil
}

func checkPr(db *sql.DB, userId int, exerciseId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	args := []interface{}{userId, exerciseId}
	var pr int
	err := db.QueryRowContext(ctx, `SELECT pr FROM exercise_prs WHERE (user_id, exercise_id) = ($1, $2);`, args).Scan(&pr)
	if err != nil {
		if pr == 0 {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (p PrModel) Update(pr Pr) error {
	return p.runQuery(updatePrQuery, []interface{}{pr.PersonalRecord, pr.UserId, pr.ExerciseId})
}

func (p PrModel) Delete(pr Pr) error {
	return p.runQuery(deletePrQuery, []interface{}{pr.UserId, pr.ExerciseId})
}

func (p PrModel) runQuery(query string, args []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := p.db.ExecContext(ctx, query, args...)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error while getting the value of rows affected")
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("No rows were affected! Please check logs!")
		return ErrRecordNotFound
	}

	if err != nil {
		return err
	}
	return nil
}

func (p PrModel) GetAll(userId int) ([]ConsolidatedPr, error) {
	var prList []ConsolidatedPr

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	args := []interface{}{userId}

	rows, err := p.db.QueryContext(ctx, selectPrByUserId, args...)
	if err != nil {
		fmt.Printf("error while fetching rows with user id: %d", userId)
		return nil, err
	}

	for rows.Next() {
		// user_id, exercise_prs.exercise_id, exercise_name, exercise_description, pr
		var pr ConsolidatedPr
		err = rows.Scan(
			&pr.UserId,
			&pr.ExerciseId,
			&pr.ExerciseName,
			&pr.ExerciseDescription,
			&pr.PersonalRecord)
		if err != nil {
			fmt.Printf("error while scanning row with user id: %d", userId)
			return nil, err
		}

		prList = append(prList, pr)
	}
	return prList, nil
}

func (p PrModel) Get(userId int, exerciseId int) (*ConsolidatedPr, error) {
	pr := ConsolidatedPr{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	args := []interface{}{userId, exerciseId}
	query := selectPrQueryByBoth

	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&pr.UserId,
		&pr.ExerciseId,
		&pr.ExerciseName,
		&pr.ExerciseDescription,
		&pr.PersonalRecord)

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

func ValidatePr(v *validator.Validator, pr *Pr, prRequired bool) {
	v.Check(pr.UserId >= 1, "User id", "must be >= 1")
	v.Check(pr.ExerciseId >= 1, "Exercise id", "must be >= 1")
	if prRequired {
		v.Check(pr.PersonalRecord > 0, "Personal Record", "must be > 0")
	}
}
