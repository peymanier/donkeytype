-- +goose up
CREATE TABLE options
(
    id text PRIMARY KEY,
    choice_id text NOT NULL,
    value     text NOT NULL
);

-- +goose down
DROP TABLE options;