CREATE TABLE IF NOT EXISTS exercise_prs (
    user_id int,
    exercise_id int REFERENCES exercises(exercise_id),
    pr int,
    PRIMARY KEY (user_id, exercise_id)
);
