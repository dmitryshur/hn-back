CREATE TABLE IF NOT EXISTS comments
(
    id       bigserial PRIMARY KEY,
    deleted  bool    NOT NULL DEFAULT false,
    type     text    NOT NULL,
    by       text,
    time     timestamp(0) with time zone,
    dead     bool    NOT NULL DEFAULT false,
    kids     integer[],
    parent   integer,
    story_id integer NOT NULL REFERENCES stories ON DELETE CASCADE,
    text     text
)
