BEGIN;

ALTER TABLE user_calorie_logs
ADD COLUMN log_status CHAR(1) CHECK (log_status IN ('P', 'D')); -- P = pending, D = Done.

CREATE TABLE IF NOT EXISTS user_caloric_balance (
  id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  calorie_log_id UUID NOT NULL,
  caloric_balance DECIMAL(6, 2) NOT NULL
);

ALTER TABLE user_calorie_logs
ALTER COLUMN log_status SET DEFAULT 'P';

END;