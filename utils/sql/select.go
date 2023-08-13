// WARNING: This is a proof of concept and isn't production ready.

package sql

import (
	"context"

	pg "movies/utils/pg"

	goqu "github.com/doug-martin/goqu/v9"
	exp "github.com/doug-martin/goqu/v9/exp"
	pgx "github.com/jackc/pgx/v4"
	lo "github.com/samber/lo"
)

func Read[M Table]() *readQuery[M] {
	return &readQuery[M]{dataset: pg.SQLBuilder().From((*new(M)).TableName())}
}

type readQuery[M any] struct {
	empty   bool
	dataset *goqu.SelectDataset
}

func (d *readQuery[M]) append(clause *goqu.SelectDataset) *readQuery[M] {
	d.dataset = clause
	return d
}

/*============================================================================*/
/*=====*                            Executor                            *=====*/
/*============================================================================*/

func (d readQuery[M]) Count(ctx context.Context, tx pg.Tx) (int, error) {
	if d.empty {
		return 0, nil
	}

	sql, args, err := d.dataset.ToSQL()
	if err != nil {
		return 0, err
	}
	var count int
	return count, pg.Get(ctx, tx, &count, sql, args...)
}

func (d readQuery[M]) FindOne(ctx context.Context, tx pg.Tx) (*M, error) {
	if d.empty {
		return nil, pgx.ErrNoRows
	}

	sql, args, err := d.dataset.ToSQL()
	if err != nil {
		return nil, err
	}
	item := new(M)
	return item, pg.Get(ctx, tx, item, sql, args...)
}

func (d readQuery[M]) FindAll(ctx context.Context, tx pg.Tx) ([]*M, error) {
	if d.empty {
		return nil, nil
	}

	sql, args, err := d.dataset.ToSQL()
	if err != nil {
		return nil, err
	}
	var items []*M
	return items, pg.Select(ctx, tx, &items, sql, args...)
}

func (d readQuery[M]) Get(ctx context.Context, tx pg.Tx, dst any) error {
	if d.empty {
		return pgx.ErrNoRows
	}

	sql, args, err := d.dataset.ToSQL()
	if err != nil {
		return err
	}
	return pg.Get(ctx, tx, dst, sql, args...)
}

func (d readQuery[M]) Sel(ctx context.Context, tx pg.Tx, dst any) error {
	if d.empty {
		return pgx.ErrNoRows
	}

	sql, args, err := d.dataset.ToSQL()
	if err != nil {
		return err
	}
	return pg.Select(ctx, tx, dst, sql, args...)
}

/*============================================================================*/
/*=====*                             Clause                             *=====*/
/*============================================================================*/

func (d *readQuery[M]) Check(values ...any) *readQuery[M] {
	d.empty = lo.Ternary(d.empty, true, len(values) == 0)
	return d
}

func (d *readQuery[M]) Where(expressions ...exp.Expression) *readQuery[M] {
	return d.append(d.dataset.Where(expressions...))
}

func (d *readQuery[M]) Order(order ...exp.OrderedExpression) *readQuery[M] {
	return d.append(d.dataset.Order(order...))
}

func (d *readQuery[M]) Limit(limit uint) *readQuery[M] {
	return d.append(d.dataset.Limit(limit))
}

func (d *readQuery[M]) GroupBy(groupBy ...any) *readQuery[M] {
	return d.append(d.dataset.GroupBy(groupBy...))
}

func (d *readQuery[M]) Select(selects ...any) *readQuery[M] {
	return d.append(d.dataset.Select(selects...))
}

func (d *readQuery[M]) SelectI(col string) *readQuery[M] {
	return d.append(d.dataset.Select(I(col)))
}

func (d *readQuery[M]) SelectDistinct(col string) *readQuery[M] {
	return d.append(d.dataset.Select(Distinct(col)))
}

func (d *readQuery[M]) From(table ...any) *readQuery[M] {
	return d.append(d.dataset.From(table...))
}

func (d *readQuery[M]) FromModel(model Table, as string) *readQuery[M] {
	return d.append(d.dataset.From(TAs(model, as)))
}

func (d *readQuery[M]) Join(table exp.Expression, condition exp.JoinCondition) *readQuery[M] {
	return d.append(d.dataset.Join(table, condition))
}

func (d *readQuery[M]) LeftJoin(table exp.Expression, condition exp.JoinCondition) *readQuery[M] {
	return d.append(d.dataset.LeftJoin(table, condition))
}

func (d *readQuery[M]) FullJoin(table exp.Expression, condition exp.JoinCondition) *readQuery[M] {
	return d.append(d.dataset.FullJoin(table, condition))
}

func (d *readQuery[M]) UnionAll(other *goqu.SelectDataset) *readQuery[M] {
	return d.append(d.dataset.UnionAll(other))
}

func (d *readQuery[M]) Having(expressions ...exp.Expression) *readQuery[M] {
	return d.append(d.dataset.Having(expressions...))
}

func (d *readQuery[M]) Raw() *goqu.SelectDataset {
	return d.dataset
}
