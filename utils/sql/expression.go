package sql

import (
	goqu "github.com/doug-martin/goqu/v9"
	exp "github.com/doug-martin/goqu/v9/exp"
)

/*============================================================================*/
/*=====*                              Func                              *=====*/
/*============================================================================*/

type Record = goqu.Record

var I = goqu.I
var T = goqu.T
var V = goqu.V
var L = goqu.L
var On = goqu.On
var MAX = goqu.MAX
var Or = goqu.Or
var And = goqu.And
var Func = goqu.Func
var Count = goqu.COUNT
var COALESCE = goqu.COALESCE
var Cast = goqu.Cast
var Distinct = func(column string) exp.SQLFunctionExpression {
	return goqu.DISTINCT(I(column))
}
var TAs = func(model Table, as string) exp.AliasedExpression {
	return T(model.TableName()).As(as)
}

var SUM = func(column string, empty any) exp.SQLFunctionExpression {
	return goqu.COALESCE(goqu.SUM(column), empty)
}

var JsonKey = func(column string, key string) exp.LiteralExpression {
	return goqu.L("(? ->> ?)", goqu.I(column), key)
}

/*============================================================================*/
/*=====*                            Constant                            *=====*/
/*============================================================================*/

var NOW = Func("NOW")
var CountALL = Count("*")
