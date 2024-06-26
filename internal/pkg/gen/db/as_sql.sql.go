// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: as_sql.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const getUnreportedImageList = `-- name: GetUnreportedImageList :many
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
WHERE s.is_reported = false
`

type GetUnreportedImageListRow struct {
	ScanID                int32           `json:"scan_id"`
	ID                    uuid.UUID       `json:"id"`
	Name                  string          `json:"name"`
	Url                   string          `json:"url"`
	Source                ImageSourceType `json:"source"`
	NumberOfLayers        int32           `json:"number_of_layers"`
	Architecture          string          `json:"architecture"`
	ConnectedRepositoryID uuid.NullUUID   `json:"connected_repository_id"`
	HighestSeverity       HighestSeverity `json:"highest_severity"`
	TotalFindings         int32           `json:"total_findings"`
	ScanDateTimestamp     time.Time       `json:"scan_date_timestamp"`
	CreatedDateTimestamp  time.Time       `json:"created_date_timestamp"`
}

func (q *Queries) GetUnreportedImageList(ctx context.Context) ([]GetUnreportedImageListRow, error) {
	rows, err := q.db.QueryContext(ctx, getUnreportedImageList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUnreportedImageListRow{}
	for rows.Next() {
		var i GetUnreportedImageListRow
		if err := rows.Scan(
			&i.ScanID,
			&i.ID,
			&i.Name,
			&i.Url,
			&i.Source,
			&i.NumberOfLayers,
			&i.Architecture,
			&i.ConnectedRepositoryID,
			&i.HighestSeverity,
			&i.TotalFindings,
			&i.ScanDateTimestamp,
			&i.CreatedDateTimestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnreportedRepositoryList = `-- name: GetUnreportedRepositoryList :many
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
WHERE s.is_reported = FALSE
`

type GetUnreportedRepositoryListRow struct {
	ScanID               int32                `json:"scan_id"`
	ID                   uuid.UUID            `json:"id"`
	Name                 string               `json:"name"`
	Url                  string               `json:"url"`
	Source               RepositorySourceType `json:"source"`
	LastPush             time.Time            `json:"last_push"`
	Size                 int64                `json:"size"`
	HighestSeverity      HighestSeverity      `json:"highest_severity"`
	TotalFindings        int32                `json:"total_findings"`
	ScanDateTimestamp    time.Time            `json:"scan_date_timestamp"`
	CreatedDateTimestamp time.Time            `json:"created_date_timestamp"`
}

func (q *Queries) GetUnreportedRepositoryList(ctx context.Context) ([]GetUnreportedRepositoryListRow, error) {
	rows, err := q.db.QueryContext(ctx, getUnreportedRepositoryList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUnreportedRepositoryListRow{}
	for rows.Next() {
		var i GetUnreportedRepositoryListRow
		if err := rows.Scan(
			&i.ScanID,
			&i.ID,
			&i.Name,
			&i.Url,
			&i.Source,
			&i.LastPush,
			&i.Size,
			&i.HighestSeverity,
			&i.TotalFindings,
			&i.ScanDateTimestamp,
			&i.CreatedDateTimestamp,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertInputImage = `-- name: InsertInputImage :one
INSERT INTO "input_image" (
    id,
    name,
    url,
    source,
    number_of_layers,
    architecture,
    created_date_timestamp
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING indx, id, repository_id, name, url, type, source, number_of_layers, architecture, created_date_timestamp
`

type InsertInputImageParams struct {
	ID                   uuid.UUID        `json:"id"`
	Name                 string           `json:"name"`
	Url                  string           `json:"url"`
	Source               ImageSourceType  `json:"source"`
	NumberOfLayers       int32            `json:"number_of_layers"`
	Architecture         ArchitectureType `json:"architecture"`
	CreatedDateTimestamp time.Time        `json:"created_date_timestamp"`
}

func (q *Queries) InsertInputImage(ctx context.Context, arg InsertInputImageParams) (InputImage, error) {
	row := q.db.QueryRowContext(ctx, insertInputImage,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.Source,
		arg.NumberOfLayers,
		arg.Architecture,
		arg.CreatedDateTimestamp,
	)
	var i InputImage
	err := row.Scan(
		&i.Indx,
		&i.ID,
		&i.RepositoryID,
		&i.Name,
		&i.Url,
		&i.Type,
		&i.Source,
		&i.NumberOfLayers,
		&i.Architecture,
		&i.CreatedDateTimestamp,
	)
	return i, err
}

const insertInputRepository = `-- name: InsertInputRepository :one
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
RETURNING indx, id, name, url, type, source, size, last_push, created_date_timestamp
`

type InsertInputRepositoryParams struct {
	ID                   uuid.UUID            `json:"id"`
	Name                 string               `json:"name"`
	Url                  string               `json:"url"`
	Source               RepositorySourceType `json:"source"`
	LastPush             time.Time            `json:"last_push"`
	Size                 int64                `json:"size"`
	CreatedDateTimestamp time.Time            `json:"created_date_timestamp"`
}

func (q *Queries) InsertInputRepository(ctx context.Context, arg InsertInputRepositoryParams) (InputRepository, error) {
	row := q.db.QueryRowContext(ctx, insertInputRepository,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.Source,
		arg.LastPush,
		arg.Size,
		arg.CreatedDateTimestamp,
	)
	var i InputRepository
	err := row.Scan(
		&i.Indx,
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Type,
		&i.Source,
		&i.Size,
		&i.LastPush,
		&i.CreatedDateTimestamp,
	)
	return i, err
}

const insertInputScan = `-- name: InsertInputScan :one
INSERT INTO "scan" (
    scan_id,
    resource_id,
    resource_type,
    highest_severity,
    total_findings,
    scan_date_timestamp
) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (scan_id,resource_id) DO UPDATE SET scan_id=$1
RETURNING indx, scan_id, resource_id, resource_type, highest_severity, total_findings, is_reported, scan_date_timestamp
`

type InsertInputScanParams struct {
	ScanID            int32           `json:"scan_id"`
	ResourceID        uuid.UUID       `json:"resource_id"`
	ResourceType      ResourceType    `json:"resource_type"`
	HighestSeverity   HighestSeverity `json:"highest_severity"`
	TotalFindings     int32           `json:"total_findings"`
	ScanDateTimestamp time.Time       `json:"scan_date_timestamp"`
}

func (q *Queries) InsertInputScan(ctx context.Context, arg InsertInputScanParams) (Scan, error) {
	row := q.db.QueryRowContext(ctx, insertInputScan,
		arg.ScanID,
		arg.ResourceID,
		arg.ResourceType,
		arg.HighestSeverity,
		arg.TotalFindings,
		arg.ScanDateTimestamp,
	)
	var i Scan
	err := row.Scan(
		&i.Indx,
		&i.ScanID,
		&i.ResourceID,
		&i.ResourceType,
		&i.HighestSeverity,
		&i.TotalFindings,
		&i.IsReported,
		&i.ScanDateTimestamp,
	)
	return i, err
}

const setBulkUpdateScanReportTrue = `-- name: SetBulkUpdateScanReportTrue :execrows
UPDATE scan
SET is_reported=true
WHERE scan_id = ANY ($1::integer[])
`

func (q *Queries) SetBulkUpdateScanReportTrue(ctx context.Context, dollar_1 []int32) (int64, error) {
	result, err := q.db.ExecContext(ctx, setBulkUpdateScanReportTrue, pq.Array(dollar_1))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateConnection = `-- name: UpdateConnection :exec
UPDATE "input_image" as ii SET
    repository_id = $1
FROM "input_repository" as ir
WHERE ir.id  = $1 and  ii.id = $2
`

type UpdateConnectionParams struct {
	RepositoryID uuid.NullUUID `json:"repository_id"`
	ID           uuid.UUID     `json:"id"`
}

func (q *Queries) UpdateConnection(ctx context.Context, arg UpdateConnectionParams) error {
	_, err := q.db.ExecContext(ctx, updateConnection, arg.RepositoryID, arg.ID)
	return err
}
