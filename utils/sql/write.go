package sql

import (
	"context"

	pg "movies/utils/pg"

	goqu "github.com/doug-martin/goqu/v9"
	exp "github.com/doug-martin/goqu/v9/exp"
	pgtype "github.com/jackc/pgtype"
)

func Create[
	M interface{ TableName() string },
](ctx context.Context, tx pg.Tx, data M, record Record) error {
	sql, args, err := pg.SQLBuilder().
		Insert(data.TableName()).
		Rows(record).
		Returning("*").
		ToSQL()
	if err != nil {
		return err
	}
	return pg.Get(ctx, tx, data, sql, args...)
}

func Update[
	M interface{ TableName() string },
](ctx context.Context, tx pg.Tx, data M, update bool, record Record, expressions exp.Expression) error {
	if update {
		record["updated_at"] = NOW
	}

	sql, args, err := pg.SQLBuilder().
		Update(data.TableName()).
		Set(record).
		Where(expressions).
		Returning("*").
		ToSQL()
	if err != nil {
		return err
	}
	return pg.Get(ctx, tx, data, sql, args...)
}

func UpdateByPK[
	M interface {
		GetPK() pgtype.UUID
		TableName() string
	},
](ctx context.Context, tx pg.Tx, data M, update bool, record Record) error {
	return Update(ctx, tx, data, update, record, goqu.And(
		goqu.I("id").Eq(data.GetPK()),
		goqu.I("deleted_at").IsNull(),
	))
}

func SoftDeleteByPK[
	M interface {
		GetPK() pgtype.UUID
		TableName() string
	},
](ctx context.Context, tx pg.Tx, data M) error {
	return Update(ctx, tx, data, false,
		Record{"deleted_at": "NOW"},
		goqu.And(
			goqu.I("id").Eq(data.GetPK()),
			goqu.I("deleted_at").IsNull(),
		))
}

//// ----> V2

func Updatev2[
	M interface{ TableName() string },
](ctx context.Context, tx pg.Tx, data M, record Record, expressions exp.Expression) error {
	sql, args, err := pg.SQLBuilder().
		Update(data.TableName()).
		Set(record).
		Where(expressions).
		Returning("*").
		ToSQL()
	if err != nil {
		return err
	}
	return pg.Get(ctx, tx, data, sql, args...)
}

func UpdateByPKv2[
	M interface {
		TableName() string
		GetPK() pgtype.UUID
		GetUpdated() Updated
	},
](ctx context.Context, tx pg.Tx, data M, track bool, record Record) error {
	if track {
		record["updated_at"] = NOW
		record["updated_by"] = data.GetUpdated().UpdatedBy
	}

	return Updatev2(ctx, tx, data, record, goqu.And(
		goqu.I("id").Eq(data.GetPK()),
		goqu.I("deleted_at").IsNull(),
	))
}

func ArchiveByPKv2[
	M interface {
		TableName() string
		GetPK() pgtype.UUID
		GetArchived() Archived
	},
](ctx context.Context, tx pg.Tx, data M) error {
	return Updatev2(ctx, tx, data,
		Record{
			"archived_at": NOW,
			"archived_by": data.GetArchived().ArchivedBy,
		},
		goqu.And(
			goqu.I("id").Eq(data.GetPK()),
			goqu.I("deleted_at").IsNull(),
		))
}

func SoftDeleteByPKv2[
	M interface {
		GetPK() pgtype.UUID
		TableName() string
		GetDeleted() Deleted
	},
](ctx context.Context, tx pg.Tx, data M) error {
	return Updatev2(ctx, tx, data,
		Record{
			"deleted_at": NOW,
			"deleted_by": data.GetDeleted().DeletedBy,
		},
		goqu.And(
			goqu.I("id").Eq(data.GetPK()),
			goqu.I("deleted_at").IsNull(),
		))
}

func HardDelete[
	M interface{ TableName() string },
](ctx context.Context, tx pg.Tx, data M, expressions exp.Expression) (int64, error) {
	sql, args, err := pg.SQLBuilder().
		Delete(data.TableName()).
		Where(expressions).
		ToSQL()
	if err != nil {
		return 0, err
	}
	ct, err := pg.Client(tx).Exec(ctx, sql, args...)
	return ct.RowsAffected(), err
}
