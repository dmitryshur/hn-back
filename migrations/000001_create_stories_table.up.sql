CREATE TABLE IF NOT EXISTS stories
(
    id          bigserial PRIMARY KEY,
    deleted     bool NOT NULL DEFAULT false,
    type        text NOT NULL,
    by          text,
    time        timestamp(0) with time zone,
    dead        bool NOT NULL DEFAULT false,
    kids        integer[],
    descendants integer,
    score       integer,
    title       text,
    url         text
)
