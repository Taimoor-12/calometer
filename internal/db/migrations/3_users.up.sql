ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);

CREATE TABLE IF NOT EXISTS user_body_details (
  id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  u_id UUID UNIQUE NOT NULL,
  age INT NOT NULL,
  height_cm INT NOT NULL,
  weight_kg DECIMAL(5,2) NOT NULL,
  gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F'))
);