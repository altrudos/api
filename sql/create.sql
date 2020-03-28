BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS drives (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  uri TEXT NOT NULL,
  amount NUMERIC NOT NULL DEFAULT 0,
  source_url TEXT,
  reddit_comment_id BIGINT,
  reddit_username TEXT,
  reddit_subreddit TEXT,
  reddit_markdown TEXT
);
CREATE UNIQUE INDEX drives_uri ON drives (uri);

CREATE TABLE IF NOT EXISTS charities (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT,
  logo_url TEXT,
  description TEXT,
  summary TEXT,
  jg_charity_id BIGINT
);
CREATE UNIQUE INDEX charities_jg_charity_id_unique ON charities (jg_charity_id);

CREATE TYPE donation_status AS ENUM ('Accepted', 'Pending', 'Rejected');

CREATE TABLE IF NOT EXISTS donations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  charity_id UUID REFERENCES charities(id),
  created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  drive_id UUID REFERENCES drives(id),
  donor_amount NUMERIC,
  donor_currency_code TEXT,
  donor_name TEXT,
  final_amount NUMERIC NOT NULL DEFAULT 0,
  last_checked TIMESTAMPTZ,
  next_check TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  message TEXT,
  message_visible BOOLEAN NOT NULL DEFAULT FALSE,
  status donation_status NOT NULL DEFAULT 'Pending',
  reference_code TEXT NOT NULL UNIQUE
);
COMMIT;
