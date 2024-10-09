BEGIN;

ALTER TABLE user_calorie_logs
ALTER COLUMN calories_burnt SET DEFAULT 0.00;

ALTER TABLE user_calorie_logs
ALTER COLUMN calories_consumed SET DEFAULT 0.00;

END;