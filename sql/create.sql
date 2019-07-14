CREATE TABLE IF NOT EXISTS drives (
  id SERIAL PRIMARY KEY,
  uri TEXT NOT NULL UNIQUE,
  amount NUMERIC NOT NULL DEFAULT 0,
  source_url TEXT,
  reddit_comment_id BIGINT,
  reddit_username TEXT,
  reddit_subreddit TEXT,
  reddit_markdown TEXT
);

CREATE TABLE IF NOT EXISTS charities (
  id SERIAL PRIMARY KEY,
  name TEXT,
  logo_url TEXT,
  description TEXT,
  summary TEXT,
  jg_charity_id BIGINT
);

CREATE TABLE IF NOT EXISTS donations (
  id SERIAL PRIMARY KEY,
  drive_id BIGINT REFERENCES drives(id),
  charity_id BIGINT REFERENCES charities(id),
  last_checked TIMESTAMPTZ,
  created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  reference_code TEXT NOT NULL UNIQUE,
  currency_code TEXT,
  amount NUMERIC,
  local_amount NUMERIC,
  local_currency_code TEXT,
  donor_name TEXT,
  message TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  message_visible BOOLEAN NOT NULL DEFAULT FALSE
);
