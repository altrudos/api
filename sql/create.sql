BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS drives
(
    id                UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    uri               TEXT        NOT NULL,
    amount            NUMERIC     NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source_url        TEXT NOT NULL,
    source_type       TEXT NOT NULL DEFAULT 'link',
    source_meta       JSONB NOT NULL DEFAULT '{}'
);
CREATE UNIQUE INDEX drives_uri ON drives (uri);

CREATE TABLE IF NOT EXISTS charities
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    country_code TEXT NOT NULL DEFAULT '',
    description   TEXT NOT NULL    DEFAULT '',
    feature_score INT  NOT NULL    DEFAULT 0,
    jg_charity_id BIGINT,
    logo_url      TEXT NOT NULL    DEFAULT '',
    name          TEXT NOT NULL    DEFAULT '',
    subtext       TEXT NOT NULL DEFAULT '',
    summary       TEXT NOT NULL    DEFAULT '',
    website_url   TEXT NOT NULL    DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX charities_jg_charity_id_unique ON charities (jg_charity_id);

CREATE TYPE donation_status AS ENUM ('Accepted', 'Pending', 'Rejected');

CREATE TABLE IF NOT EXISTS donations
(
    id                  UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    charity_id          UUID REFERENCES charities (id),
    created_at             TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
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
    reference_code      TEXT            NOT NULL UNIQUE,
    usd_amount INT NOT NULL DEFAULT 0
);
CREATE INDEX donation_charity_id ON donations (charity_id);
CREATE INDEX donation_drive_id ON donations (drive_id);
CREATE INDEX donation_donor_amount ON donations (donor_amount);
CREATE INDEX donation_final_amount ON donations (final_amount);
CREATE INDEX donation_usd_amount ON donations (usd_amount);
CREATE INDEX donation_created_at ON donations (created_at);

CREATE TABLE IF NOT EXISTS search_cache
(
    term    TEXT UNIQUE NOT NULL DEFAULT '',
    expires TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ids     BIGINT[]
);

COMMIT;
