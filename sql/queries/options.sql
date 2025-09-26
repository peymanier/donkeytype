-- name: ListOptions :many
SELECT *
FROM options;

-- name: AddOption :exec
INSERT INTO options (id, choice_id, value)
VALUES (?, ?, ?)
ON CONFLICT (id) DO UPDATE SET value     = excluded.value,
                               choice_id = excluded.choice_id;
