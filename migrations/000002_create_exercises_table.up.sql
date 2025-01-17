CREATE TABLE IF NOT EXISTS exercises (
    exercise_id bigserial PRIMARY KEY,
    exercise_name text NOT NULL,
    exercise_description text
);