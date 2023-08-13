package pg

import (
	goqu "github.com/doug-martin/goqu/v9"
)

// Model ...
type Model interface {
	TableName() string
}

//nolint:gochecknoinits // there is no other way to force goqu to use prepared statements
func init() {
	goqu.SetDefaultPrepared(true)
}

// SQLBuilder ...
func SQLBuilder() goqu.DialectWrapper {
	return goqu.Dialect("postgres")
}

func QueryModel(model Model, cols ...any) *goqu.SelectDataset {
	tmp := SQLBuilder().From(model.TableName())
	if len(cols) > 0 {
		return tmp.Select(cols...)
	}
	return tmp.Select(model)
}

// QuerySelect ...
func QuerySelect(model Model, where goqu.Expression) (string, []any) {
	sql, args, err := QueryModel(model).Where(where).ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// QueryUpdate ...
func QueryUpdate(model Model, values any, where goqu.Expression) (string, []any) {
	sql, args, err := SQLBuilder().Update(model.TableName()).Set(values).Where(where).Returning("*").ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// QueryCount ...
func QueryCount(model Model, where goqu.Expression) (string, []any) {
	sql, args, err := QueryModel(model, goqu.COUNT("*")).Where(where).ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// QueryInsert ...
func QueryInsert(model Model, rows ...any) (string, []any) {
	sql, args, err := SQLBuilder().Insert(model.TableName()).Rows(rows...).Returning("*").ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// QueryDelete ...
func QueryDelete(model Model, where goqu.Expression) (string, []any) {
	sql, args, err := SQLBuilder().Delete(model.TableName()).Where(where).ToSQL()
	if err != nil {
		panic(err)
	}
	return sql, args
}
