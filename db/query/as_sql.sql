-- name: InsertInputRepository :one
INSERT INTO "input_repository" (
    id,
    name,
    url,
    source,
    last_push,
    size,
    created_date_timestamp
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) where last_push < $5::timestamp
DO UPDATE SET last_push = $5, size = $6
RETURNING *;


-- name: InsertInputImage :one
INSERT INTO "input_image" (
    id,
    name,
    url,
    source,
    number_of_layers,
    architecture,
    created_date_timestamp
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateConnection :exec
UPDATE "input_image" as ii SET
    repository_id = $1
FROM "input_repository" as ir
WHERE ir.id  = $1 and  ii.id = $2;

-- name: InsertInputScan :one
INSERT INTO "scan" (
    scan_id,
    resource_id,
    resource_type,
    highest_severity,
    total_findings,
    scan_date_timestamp
) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (scan_id,resource_id) DO UPDATE SET scan_id=$1
RETURNING *;

-- name: GetUnreportedRepositoryList :many
SELECT
    s.scan_id,
    ir.id,
    ir."name",
    ir.url,
    ir."source",
    ir.last_push,
    ir."size",
    s.highest_severity,
    s.total_findings,
    s.scan_date_timestamp,
    ir.created_date_timestamp
FROM scan s
JOIN input_repository ir ON s.resource_id = ir.id
WHERE s.is_reported = FALSE;

-- name: GetUnreportedImageList :many
SELECT
    s.scan_id,
    ii.id,
    ii."name",
    ii.url,
    ii."source",
    ii.number_of_layers ,
    ii.architecture::text as architecture,
    ii.repository_id as "connected_repository_id",
    s.highest_severity,
    s.total_findings,
    s.scan_date_timestamp,
    ii.created_date_timestamp
FROM scan s
JOIN input_image ii ON s.resource_id = ii.id
WHERE s.is_reported = false;


-- name: SetBulkUpdateScanReportTrue :execrows
UPDATE scan
SET is_reported=true
WHERE scan_id = ANY ($1::integer[]);