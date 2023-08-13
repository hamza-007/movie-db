package pg

import (
	"context"
	"sync"

	logger "movies/utils/logger"

	goqu "github.com/doug-martin/goqu/v9"
	dbscan "github.com/georgysavva/scany/dbscan"
	pgxscan "github.com/georgysavva/scany/pgxscan"
	pgtype "github.com/jackc/pgtype"
	viper "github.com/spf13/viper"
)

var (
	scanAPI  *pgxscan.API
	scanOnce sync.Once
)

func scany() *pgxscan.API {
	scanOnce.Do(func() {
		config, err := pgxscan.NewDBScanAPI(
			dbscan.WithAllowUnknownColumns(true),
		)
		if err != nil {
			panic(err)
		}
		scanAPI, err = pgxscan.NewAPI(config)
		if err != nil {
			panic(err)
		}
	})
	return scanAPI
}

func Get(ctx context.Context, tx Tx, dst any, query string, args ...any) error {
	err := scany().Get(ctx, Client(tx), dst, query, args...)
	if err != nil {
		if viper.GetBool("log-db") {
			logger.AddCtxLabel(ctx, "sql", query)
		}
	}
	return err
}

func Select(ctx context.Context, tx Tx, dst any, query string, args ...any) error {
	err := scany().Select(ctx, Client(tx), dst, query, args...)
	if err != nil {
		if viper.GetBool("log-db") {
			logger.AddCtxLabel(ctx, "sql", query)
		}
	}
	return err
}

func Create[
	M interface {
		GetPK() pgtype.UUID
		TableName() string
	},
](ctx context.Context, tx Tx, data M, rows ...any) error {
	sql, args := QueryInsert(data, rows)
	return Get(ctx, tx, data, sql, args...)
}

func Update[
	M interface {
		GetPK() pgtype.UUID
		TableName() string
	},
](ctx context.Context, tx Tx, data M, record goqu.Record) error {
	sql, args := QueryUpdate(data, record, goqu.And(
		goqu.I("id").Eq(data.GetPK()),
		goqu.I("deleted_at").IsNull(),
	))
	return Get(ctx, tx, data, sql, args...)
}
