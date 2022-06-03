CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS announcement (
    id text PRIMARY KEY,
    "text" text NOT NULL,
    image_url text NOT NULL,
    position geography(point) NOT NULL,
    created_at timestamp NOT NULL
);

CREATE INDEX announcement_position_idx ON announcement USING GIST (position, created_at);