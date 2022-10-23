CREATE TABLE IF NOT EXISTS stories
(
    id          integer PRIMARY KEY,
    deleted     bool,
    type        text NOT NULL,
    by          text,
    time        integer,
    dead        bool,
    kids        integer[],
    descendants integer,
    score       integer,
    title       text,
    url         text
)
