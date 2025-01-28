CREATE OR REPLACE FUNCTION pr_updating_function()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
DECLARE
    max_pr int;
    local_max_pr int;
    current_max_pr int;
BEGIN
    SELECT MAX(nums) INTO local_max_pr FROM unnest(new.weights) AS nums;
    SELECT pr INTO current_max_pr FROM exercise_prs WHERE (user_id, exercise_id) = (new.user_id, new.exercise_id);

    -- If the value of current_max_pr is NULL, then we know that it is a new pr.
    IF current_max_pr IS NULL THEN
        INSERT INTO exercise_prs(USER_ID, EXERCISE_ID, PR) VALUES(new.user_id, new.exercise_id, local_max_pr);
    END IF;
    -- If the value for this already exists then we need to update the value accordingly
    -- If local_max_pr > current_max_pr then we update the exercise_prs table for the (user,exercise)
    -- Otherwise we do nothing as it already updated
    IF local_max_pr > current_max_pr THEN
        UPDATE exercise_prs
        SET
            pr = local_max_pr
        WHERE (user_id, exercise_id) = (new.user_id, new.exercise_id);
    END IF;

    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER pr_updating_trigger
    BEFORE INSERT OR UPDATE
    ON workouts_table
    FOR EACH ROW
EXECUTE PROCEDURE pr_updating_function();

-- TRIGGER FOR DELETING FROM PR TABLE WHEN THERE IS A MAX RECORD WORKOUT DELETED FROM THE WORKOUT TABLE

CREATE OR REPLACE TRIGGER pr_updating_trigger
    BEFORE INSERT OR UPDATE
    ON workouts_table
    FOR EACH ROW
EXECUTE PROCEDURE pr_updating_function();

-- TRIGGER FOR DELETING FROM PR TABLE WHEN THERE IS A MAX RECORD WORKOUT DELETED FROM THE WORKOUT TABLE

CREATE OR REPLACE TRIGGER pr_updating_trigger
    BEFORE INSERT OR UPDATE
    ON workouts_table
    FOR EACH ROW
EXECUTE PROCEDURE pr_updating_function();

-- TRIGGER FOR DELETING FROM PR TABLE WHEN THERE IS A MAX RECORD WORKOUT DELETED FROM THE WORKOUT TABLE

CREATE OR REPLACE FUNCTION pr_deleting_function()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
DECLARE
    new_max_pr int;
BEGIN
    -- I have the user id and the exercise id, I will recompute the PR for this exercise
    -- after deleting it from the prs table
    DELETE FROM exercise_prs WHERE (exercise_prs.user_id, exercise_prs.exercise_id) = (old.user_id, old.exercise_id);
    -- now recompute the new max pr
    SELECT MAX(nums)
    INTO new_max_pr
    FROM (
             SELECT unnest(weights) AS nums
             FROM workouts_table
         ) wtn;

    INSERT INTO exercise_prs(USER_ID, EXERCISE_ID, PR) VALUES(old.user_id, old.exercise_id, new_max_pr);
    RETURN NULL;
END;
$$;


CREATE OR REPLACE TRIGGER pr_deletion_trigger
    AFTER DELETE
    ON workouts_table
    FOR EACH ROW
EXECUTE PROCEDURE pr_deleting_function();