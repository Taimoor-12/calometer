BEGIN;

ALTER TABLE user_caloric_balance ADD CONSTRAINT unique_calorie_log_id UNIQUE (calorie_log_id);

END;