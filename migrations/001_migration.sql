CREATE SCHEMA IF NOT EXISTS shortener;

CREATE TABLE IF NOT EXISTS shortener.users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid()
);

CREATE TABLE IF NOT EXISTS shortener.shortener(
    correlation_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    short_url varchar(55) NOT NULL,
    original_url varchar(255) NOT NULL,
    is_deleted bool NOT NULL DEFAULT false,
    user_id uuid NOT NULL REFERENCES shortener.users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx on shortener.shortener(original_url);
