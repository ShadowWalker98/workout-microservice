CREATE TABLE IF NOT EXISTS workouts_table (
    workout_id bigserial,
    user_id int,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    exercise_id bigserial REFERENCES exercises(exercise_id),
    duration int,
    sets int,
    reps int[],
    weights int[],
    PRIMARY KEY (workout_id, user_id)
);

ALTER TABLE workouts_table
ADD CONSTRAINT SETS_CONSTRAINTS CHECK (sets > 0);

ALTER TABLE workouts_table ADD CONSTRAINT REPS_CONSTRAINTS CHECK (array_length(reps, 1) > 0);

ALTER TABLE workouts_table ADD CONSTRAINT WEIGHTS_CONSTRAINTS CHECK (array_length(weights, 1) > 0);

ALTER TABLE workouts_table ADD CONSTRAINT LENGTH_CONSTRAINTS CHECK (array_length(reps, 1) = sets AND array_length(weights, 1) = sets);