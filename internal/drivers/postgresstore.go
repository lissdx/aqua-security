package drivers

import (
	"context"
	"database/sql"
	sqlcDb "github.com/lissdx/aqua-security/internal/pkg/gen/db"
	asLoger "github.com/lissdx/aqua-security/internal/pkg/logger"
)

type Store interface {
	sqlcDb.Querier
	GetDb() *sql.DB
	WithTx(tx *sql.Tx) *sqlcDb.Queries
	PrepareTx() (*sql.Tx, *sqlcDb.Queries, func(*sql.Tx, *error), error)
}

type SQLStore struct {
	logger asLoger.TOLogger
	db     *sql.DB
	*sqlcDb.Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB, logger asLoger.TOLogger) Store {
	return &SQLStore{
		db:      db,
		logger:  logger,
		Queries: sqlcDb.New(db),
	}
}

func (s *SQLStore) GetDb() *sql.DB {
	return s.db
}

func (s *SQLStore) WithTx(tx *sql.Tx) *sqlcDb.Queries {
	return s.Queries.WithTx(tx)
}

func (s *SQLStore) PrepareTx() (*sql.Tx, *sqlcDb.Queries, func(*sql.Tx, *error), error) {
	tx, err := s.GetDb().BeginTx(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	fn := func(tx *sql.Tx, err *error) {
		if *err != nil {
			s.logger.Debug("call prepareTx defer err: %s", (*err).Error())
			rErr := tx.Rollback()
			if rErr != nil {
				s.logger.Error(rErr.Error())
			}
		}
	}

	qTx := s.WithTx(tx)

	return tx, qTx, fn, nil
}
