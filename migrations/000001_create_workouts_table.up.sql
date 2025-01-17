CREATE TABLE IF NOT EXISTS workouts_table (
    workout_id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    exercise_id bigserial REFERENCES exercises(exercise_id),
    duration int,
    sets int,
    reps int[],
    weights int[]
);

ALTER TABLE workouts_table
ADD CONSTRAINT SETS_CONSTRAINTS CHECK (sets > 0);

ALTER TABLE workouts_table ADD CONSTRAINT REPS_CONSTRAINTS CHECK (length(reps) > 0);

ALTER TABLE workouts_table ADD CONSTRAINT WEIGHTS_CONSTRAINTS CHECK (length(weights) > 0);