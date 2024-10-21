-- name: ListExcludedOrigins :many
SELECT
  origin
FROM
  moz_perms
WHERE
  type = 'cookie'
  AND permission = 1
  AND expireTime = 0;

-- name: DeleteAll :execrows
DELETE FROM moz_perms
WHERE
  type = 'cookie'
  AND permission = 1
  AND expireTime = 0;

-- name: CountExcludedDomain :one
SELECT
  COUNT(*)
FROM
  moz_perms
WHERE
  type = 'cookie'
  AND permission = 1
  AND expireTime = 0
  AND (
    origin = concat('https://', cast(sqlc.arg(domain) as text))
    OR origin = concat('http://', cast(sqlc.arg(domain) as text))
  );

-- name: CreateExcludedDomain :exec
INSERT INTO
  moz_perms (
    origin,
    type,
    permission,
    expireType,
    expireTime,
    modificationTime
  )
VALUES
  (concat('https://', cast(sqlc.arg(domain) as text)), 'cookie', 1, 0, 0, cast(sqlc.arg(now) as integer)),
  (concat('http://', cast(sqlc.arg(domain) as text)), 'cookie', 1, 0, 0, cast(sqlc.arg(now) as integer));
