CREATE TABLE IF NOT EXISTS comments
(
    id       integer PRIMARY KEY,
    deleted  bool,
    type     text    NOT NULL,
    by       text,
    time     integer,
    dead     bool,
    kids     integer[],
    parent   integer,
    story_id integer NOT NULL REFERENCES stories ON DELETE CASCADE,
    text     text
)
