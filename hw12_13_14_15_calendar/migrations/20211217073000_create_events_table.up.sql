CREATE TABLE IF NOT EXISTS events (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	title TEXT NOT NULL,
	date TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
	duration INTEGER NOT NULL,
	description TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
