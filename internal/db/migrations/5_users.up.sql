BEGIN;

CREATE TABLE IF NOT EXISTS user_weight_goal (
  id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  u_id UUID UNIQUE NOT NULL,
  goal CHAR(1) NOT NULL CHECK (goal IN ('G', 'L')) -- G == Gain, L == Lose
);

END;