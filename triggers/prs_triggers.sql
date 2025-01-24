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

CREATE OR REPLACE FUNCTION pr_deleting_function()
    RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
DECLARE
    local_exercise_id int;
    local_user_id int;
    local_weights int[];
    local_pr int;
    new_max_pr int;
    existing_pr int;
BEGIN

    -- I am getting the user id and exercise id from the workout info
    SELECT user_id, exercise_id, weights
    INTO local_user_id, local_exercise_id, local_pr
    FROM workouts_table
    WHERE (workout_id) = (old.workout_id);

    -- getting the max weight in this workout
    SELECT MAX(nums) INTO local_pr FROM unnest(local_weights) AS nums;

    -- checking if the max weight in the deleted workout is the max pr of the person

    SELECT pr INTO existing_pr FROM exercise_prs
    WHERE (exercise_prs.user_id, exercise_prs.exercise_id) = (local_user_id, local_exercise_id);

    -- if they are equal then we have to update the exercise prs table
    -- otherwise the deletion does not impact that table at all

    IF local_pr = existing_pr THEN
        -- if they are equal then we replace it with the next best one
        -- putting the next best weight pr into max_pr
        SELECT MAX(weights) WHERE weights = (SELECT MAX(weights)
                                             FROM workouts_table
                                             WHERE
                                                     (user_id, exercise_id) =
                                                     (new.user_id, new.exercise_id)
                                               AND workout_id <> old.workout_id
        )

        INTO new_max_pr;
        -- if max_pr is null after this deletion then we remove the entry from the prs table
        IF new_max_pr IS NULL THEN
            DELETE FROM exercise_prs
            WHERE (exercise_prs.exercise_id, exercise_prs.user_id) = (local_exercise_id, local_user_id);
        ELSE
            -- if it is not null then we update the exercise prs table with the new_max_pr
            UPDATE exercise_prs
            SET
                pr = new_max_pr
            WHERE
                    (user_id, exercise_id) = (local_user_id, local_exercise_id);
        END IF;
    END IF;
END;
$$;


CREATE OR REPLACE TRIGGER pr_deletion_trigger
    BEFORE DELETE
    ON workouts_table
    FOR EACH STATEMENT
EXECUTE PROCEDURE pr_deleting_function();