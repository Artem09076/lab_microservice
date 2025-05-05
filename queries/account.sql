-- name: GetUserByUsermame :one
SELECT id, username, password, enabled
FROM public.account
WHERE username=$1
;

-- name: GetUserByID :one
SELECT id, username, password, enabled
FROM public.account
WHERE id=$1
;

-- name: GetActiveUsers :many
SELECT id, username, password, enabled
FROM public.account
WHERE enabled=TRUE
;