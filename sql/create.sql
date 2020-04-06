BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE source_type AS ENUM ('reddit_comment', 'reddit_post', 'url');

CREATE TABLE IF NOT EXISTS drives
(
    id                UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    uri               TEXT        NOT NULL,
    amount            NUMERIC     NOT NULL DEFAULT 0,
    created           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source_url        TEXT,
    source_type       source_type,
    source_key        TEXT NOT NULL,
    source_meta       JSONB NOT NULL DEFAULT '{}'
);
CREATE UNIQUE INDEX drives_uri ON drives (uri);
CREATE UNIQUE INDEX drives_source ON drives (source_type, source_key);

CREATE TABLE IF NOT EXISTS charities
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          TEXT NOT NULL    DEFAULT '',
    logo_url      TEXT NOT NULL    DEFAULT '',
    website_url   TEXT NOT NULL    DEFAULT '',
    description   TEXT NOT NULL    DEFAULT '',
    summary       TEXT NOT NULL    DEFAULT '',
    jg_charity_id BIGINT,
    feature_score INT  NOT NULL    DEFAULT 0
);
CREATE UNIQUE INDEX charities_jg_charity_id_unique ON charities (jg_charity_id);

CREATE TYPE donation_status AS ENUM ('Accepted', 'Pending', 'Rejected');

CREATE TABLE IF NOT EXISTS donations
(
    id                  UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    charity_id          UUID REFERENCES charities (id),
    created             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    drive_id            UUID REFERENCES drives (id),
    donor_amount        INT,
    donor_currency      TEXT,
    donor_name          TEXT,
    final_amount        INT             NOT NULL DEFAULT 0,
    final_currency      TEXT,
    last_checked        TIMESTAMPTZ,
    next_check          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    message             TEXT,
    message_visible     BOOLEAN         NOT NULL DEFAULT FALSE,
    status              donation_status NOT NULL DEFAULT 'Pending',
    reference_code      TEXT            NOT NULL UNIQUE
);
CREATE INDEX donation_charity_id ON donations (charity_id);
CREATE INDEX donation_drive_id ON donations (drive_id);
CREATE INDEX donation_donor_amount ON donations (donor_amount);
CREATE INDEX donation_final_amount ON donations (final_amount);
CREATE INDEX donation_created ON donations (created);

CREATE TABLE IF NOT EXISTS search_cache
(
    term    TEXT UNIQUE NOT NULL DEFAULT '',
    expires TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ids     BIGINT[]
);

COMMIT;
