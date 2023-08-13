package pg

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	errors "emperror.dev/errors"
	uuid "github.com/google/uuid"
	pgtype "github.com/jackc/pgtype"
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	lo "github.com/samber/lo"
	decimal "github.com/shopspring/decimal"
)

//nolint:gochecknoinits // Initialization is done in init()
func init() {
	gob.Register(pgtype.Bool{})
	gob.Register(pgtype.Date{})
	gob.Register(pgtype.Daterange{})
	gob.Register(pgtype.JSON{})
	gob.Register(pgtype.JSONB{})
	gob.Register(pgtype.UUID{})
	gob.Register(pgtype.Tstzrange{})
	gob.Register(pgtype.Timestamptz{})
	gob.Register(pgtype.TextArray{})
	gob.Register(pgtype.Point{})
	gob.Register(decimal.Decimal{})
	gob.Register(numeric.Numeric{})
}

/*============================================================================*/
/*=====*                              Bool                              *=====*/
/*============================================================================*/

func NewBool(status bool) pgtype.Bool {
	return pgtype.Bool{Status: pgtype.Present, Bool: status}
}

/*============================================================================*/
/*=====*                              Date                              *=====*/
/*============================================================================*/

// UnmarshalDate : Implements the graphql.Unmarshaler interface for `Date` GraphQL scalar and Go `pgtype.Date` struct.
func UnmarshalDate(v any) (pgtype.Date, error) {
	date := pgtype.Date{}
	if err := date.Set(v); err != nil {
		return date, errors.Errorf("Error not date")
	}
	return date, nil
}

func NewDate() pgtype.Date {
	return pgtype.Date{Status: pgtype.Present, Time: time.Now()}
}

func NewDateFromTime(time time.Time) pgtype.Date {
	return pgtype.Date{Status: pgtype.Present, Time: time}
}

// NewDateFromTimestamptz : Get a pgtype.Date from pgtype.Timestamptz
func NewDateFromTimestamptz(t pgtype.Timestamptz) pgtype.Date {
	return pgtype.Date(t)
}

func ParseDate(str string) (pgtype.Date, error) {
	date := pgtype.Date{Status: pgtype.Present}
	return date, date.Set(str)
}

func FormatDate(dr pgtype.Date) string {
	return lo.If(dr.InfinityModifier == pgtype.Infinity, "infinity").
		ElseIf(dr.InfinityModifier == pgtype.NegativeInfinity, "-infinity").
		ElseIf(dr.Status == pgtype.Null, "null").
		ElseIf(dr.Status == pgtype.Undefined, "undefined").
		Else(dr.Time.Format("2006-01-02"))
}

// CompareDate : Compare two dates, return -1 if a < b, 0 if a == b, 1 if a > b.
func CompareDate(a, b pgtype.Date) int {
	unixA := lo.If[int64](a.InfinityModifier == pgtype.Infinity, math.MaxInt64).
		ElseIf(a.InfinityModifier == pgtype.NegativeInfinity, math.MinInt64).
		ElseIf(a.Status == pgtype.Null, 0).
		ElseIf(a.Status == pgtype.Undefined, 0).
		Else(a.Time.Unix())

	unixB := lo.If[int64](b.InfinityModifier == pgtype.Infinity, math.MaxInt64).
		ElseIf(b.InfinityModifier == pgtype.NegativeInfinity, math.MinInt64).
		ElseIf(b.Status == pgtype.Null, 0).
		ElseIf(b.Status == pgtype.Undefined, 0).
		Else(b.Time.Unix())

	return lo.If(unixA == unixB, 0).ElseIf(unixA < unixB, -1).Else(1)
}

func IsDateInfinity(time pgtype.Date) bool {
	return time.InfinityModifier == pgtype.Infinity || time.InfinityModifier == pgtype.NegativeInfinity
}

func IsDateReal(time pgtype.Date) bool {
	return time.Status == pgtype.Present && !IsDateInfinity(time)
}

func EqualDates(a, b pgtype.Date) bool {
	tA := a.Time
	tB := b.Time
	return tA.Year() == tB.Year() && tA.Month() == tB.Month() && tA.Day() == tB.Day()
}

/*============================================================================*/
/*=====*                           Daterange                            *=====*/
/*============================================================================*/

// UnmarshalDaterange : Implements the graphql.Unmarshaler interface for `Daterange` GraphQL scalar and Go `pgtype.Daterange` struct.
func UnmarshalDaterange(v any) (pgtype.Daterange, error) {
	date := pgtype.Daterange{}
	u, ok := v.([]any)
	if !ok || len(u) != 2 {
		return date, errors.Errorf("Error not Daterange")
	}

	if err := date.Set(fmt.Sprintf("[%s,%s]", u[0], u[1])); err != nil {
		return date, errors.Errorf("Error not Daterange")
	}
	return date, nil
}

// NewDaterangeFromDates : Get a pgtype.Daterange from two *pgtype.Date.
func NewDaterangeFromDates(start, end *pgtype.Date) pgtype.Daterange {
	l := pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity}
	if start != nil {
		l = *start
	}

	u := pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.Infinity}
	if end != nil {
		u = *end
	}

	return pgtype.Daterange{
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Status:    pgtype.Present,
		Lower:     l,
		Upper:     u,
	}
}

// NewDaterangeFromTimes : Get a pgtype.Daterange from two *time.Time.
func NewDaterangeFromTimes(start, end *time.Time) pgtype.Daterange {
	l := pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity}
	if start != nil {
		l.InfinityModifier = pgtype.None
		l.Time = *start
	}

	u := pgtype.Date{Status: pgtype.Present, InfinityModifier: pgtype.Infinity}
	if end != nil {
		u.InfinityModifier = pgtype.None
		u.Time = *end
	}

	return pgtype.Daterange{
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Status:    pgtype.Present,
		Lower:     l,
		Upper:     u,
	}
}

// ParseDaterange : Get a pgtype.Daterange from string.
func ParseDaterange(v string) (pgtype.Daterange, error) {
	date := pgtype.Daterange{}
	return date, date.Set(v)
}

func NormalizeDaterange(period pgtype.Daterange) pgtype.Daterange {
	if period.UpperType == pgtype.Inclusive {
		period.UpperType = pgtype.Exclusive
		if period.Upper.InfinityModifier == pgtype.None {
			period.Upper.Time = period.Upper.Time.AddDate(0, 0, 1)
		}
	}
	return period
}

func FormatDaterange(dr pgtype.Daterange) string {
	return fmt.Sprintf(
		"%s%s,%s%s",
		lo.Ternary(dr.LowerType == pgtype.Inclusive, "[", "("),
		FormatDate(dr.Lower),
		FormatDate(dr.Upper),
		lo.Ternary(dr.UpperType == pgtype.Inclusive, "]", ")"),
	)
}

/*============================================================================*/
/*=====*                             JSONB                              *=====*/
/*============================================================================*/

// UnmarshalJSONB : Implements the graphql.Unmarshaler interface for `JSONB` GraphQL scalar and Go `pgtype.JSONB` struct.
func UnmarshalJSONB(v any) (pgtype.JSONB, error) {
	date := pgtype.JSONB{}
	if err := date.Set(v); err != nil {
		return date, errors.Errorf("Error not json")
	}
	return date, nil
}

// NewJSONBFromAny : Get a pgtype.JSONB from any
func NewJSONBFromAny(v any) (pgtype.JSONB, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return pgtype.JSONB{}, err
	}

	data := pgtype.JSONB{}
	if err := data.Set(b); err != nil {
		return pgtype.JSONB{}, err
	}
	return data, nil
}

func NewJSONBFromBytes(v []byte) pgtype.JSONB {
	data := pgtype.JSONB{}
	if err := data.Set(v); err != nil {
		panic(err)
	}
	return data
}

func EmptyJSONBObject() pgtype.JSONB {
	return NewJSONBFromBytes([]byte("{}"))
}

/*============================================================================*/
/*=====*                            Decimal                             *=====*/
/*============================================================================*/

// UnmarshalDecimal : Implements the graphql.Unmarshaler interface for `Decimal` GraphQL scalar and Go `numeric.Decimal` struct.
func UnmarshalDecimal(v any) (decimal.Decimal, error) {
	n, err := decimal.NewFromFormattedString(fmt.Sprintf("%v", v), nil)
	if err != nil {
		return n, errors.Errorf("Error not decimal")
	}
	return n, nil
}

/*============================================================================*/
/*=====*                            Numeric                             *=====*/
/*============================================================================*/

// UnmarshalNumeric : Implements the graphql.Unmarshaler interface for `Numeric` GraphQL scalar and Go `numeric.Numeric` struct.
func UnmarshalNumeric(v any) (numeric.Numeric, error) {
	n := numeric.Numeric{}
	if err := n.Set(v); err != nil {
		return n, errors.Errorf("Error not numeric")
	}
	return n, nil
}

func newNumeric(value decimal.Decimal) numeric.Numeric {
	return numeric.Numeric{Status: pgtype.Present, Decimal: value}
}

func NullNumeric() numeric.Numeric {
	return numeric.Numeric{Status: pgtype.Null}
}

func NewNumericFromFloat64(value float64) numeric.Numeric {
	return newNumeric(decimal.NewFromFloat(value))
}

func NewNumericFromInt64(value int64) numeric.Numeric {
	return newNumeric(decimal.NewFromInt(value))
}

func NewNumericFromDecimal(value decimal.Decimal) numeric.Numeric {
	return newNumeric(value)
}

/*============================================================================*/
/*=====*                             Point                              *=====*/
/*============================================================================*/

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// UnmarshalPoint : Implements the graphql.Unmarshaler interface for `Point` GraphQL scalar and Go `pgtype.Point` struct.
func UnmarshalPoint(v any) (pgtype.Point, error) {
	point := pgtype.Point{}
	u, ok := v.([]any)
	if !ok {
		return point, errors.Errorf("Not a Point")
	}
	if len(u) != 2 {
		point.Status = pgtype.Null
		return point, nil
	}
	if err := point.Set(fmt.Sprintf("(%v,%v)", u[0], u[1])); err != nil {
		return point, errors.Errorf("Not a point")
	}
	return point, nil
}

/*============================================================================*/
/*=====*                           TextArray                            *=====*/
/*============================================================================*/

// UnmarshalTextArray : Implements the graphql.Unmarshaler interface for `TextArray` GraphQL scalar and Go `pgtype.TextArray` struct.
func UnmarshalTextArray(v any) (pgtype.TextArray, error) {
	array := pgtype.TextArray{}
	if err := array.Set(v); err != nil {
		return array, errors.Errorf("Error not array")
	}
	return array, nil
}

func TextArrayToPrimitive(array pgtype.TextArray) []string {
	return lo.Reduce(array.Elements, func(agg []string, item pgtype.Text, _ int) []string {
		return append(agg, item.String)
	}, []string{})
}

func PrimitiveToTextArray(array []string) (pgtype.TextArray, error) {
	result := pgtype.TextArray{}
	err := result.Set(array)
	return result, err
}

/*============================================================================*/
/*=====*                          Timestamptz                           *=====*/
/*============================================================================*/

// UnmarshalTimestamptz : Implements the graphql.Unmarshaler interface for `Timestamptz` GraphQL scalar and Go `pgtype.Timestamptz` struct.
func UnmarshalTimestamptz(v any) (pgtype.Timestamptz, error) {
	t, err := time.Parse(time.RFC3339, v.(string))
	if err != nil {
		return pgtype.Timestamptz{}, errors.Errorf("Error not RFC3339")
	}
	return pgtype.Timestamptz{Time: t, Status: pgtype.Present}, nil
}

// NewTimestamptz : Get a pgtype.Timestamptz set to Now.
func NewTimestamptz() pgtype.Timestamptz {
	return pgtype.Timestamptz{Status: pgtype.Present, Time: time.Now()}
}

func NewNullTimestamptz() pgtype.Timestamptz {
	return pgtype.Timestamptz{Status: pgtype.Null}
}

func NewTimestamptzFromTime(time time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Status: pgtype.Present, Time: time}
}

// IsTstzInfinity : Check if a Timestamptz is infinity or negative infinity.
func IsTstzInfinity(time pgtype.Timestamptz) bool {
	return time.InfinityModifier == pgtype.Infinity || time.InfinityModifier == pgtype.NegativeInfinity
}

func IsTstzReal(time pgtype.Timestamptz) bool {
	return time.Status == pgtype.Present && !IsTstzInfinity(time)
}

/*============================================================================*/
/*=====*                           Tstzrange                            *=====*/
/*============================================================================*/

// UnmarshalTstzrange : Implements the graphql.Unmarshaler interface for `Tstzrange` GraphQL scalar and Go `pgtype.Tstzrange` struct.
func UnmarshalTstzrange(v any) (pgtype.Tstzrange, error) {
	date := pgtype.Tstzrange{}
	u, ok := v.([]any)
	if !ok || len(u) != 2 {
		return date, errors.Errorf("Error not Tstzrange")
	}

	if err := date.Set(fmt.Sprintf("[%s,%s)", u[0], u[1])); err != nil {
		return date, errors.Errorf("Error not Tstzrange")
	}
	return date, nil
}

// NewTstzrangeFromTstzs : Get a pgtype.Tstzrange from two *pgtype.Timestamptz.
func NewTstzrangeFromTstzs(start, end *pgtype.Timestamptz) pgtype.Tstzrange {
	l := pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity}
	if start != nil {
		l = *start
	}

	u := pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.Infinity}
	if end != nil {
		u = *end
	}

	return pgtype.Tstzrange{
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Status:    pgtype.Present,
		Lower:     l,
		Upper:     u,
	}
}

func NewTstzrangeFromTimes(start, end *time.Time) pgtype.Tstzrange {
	l := pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.NegativeInfinity}
	if start != nil {
		l.Time = *start
		l.InfinityModifier = pgtype.None
	}

	u := pgtype.Timestamptz{Status: pgtype.Present, InfinityModifier: pgtype.Infinity}
	if end != nil {
		u.Time = *end
		u.InfinityModifier = pgtype.None
	}

	return pgtype.Tstzrange{
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Status:    pgtype.Present,
		Lower:     l,
		Upper:     u,
	}
}

func TstzrangeContain(period pgtype.Tstzrange, t time.Time) bool {
	// If the period is null, it does not contain anything.
	if period.Status == pgtype.Null {
		return false
	}
	// If the period is infinite, it contains everything.
	if period.Lower.InfinityModifier == pgtype.NegativeInfinity && period.Upper.InfinityModifier == pgtype.Infinity {
		return true
	}
	// If the period is infinite on one side, it contains everything on the other side.
	if period.Lower.InfinityModifier == pgtype.NegativeInfinity {
		return t.Before(period.Upper.Time)
	}
	// If the period is infinite on one side, it contains everything on the other side.
	if period.Upper.InfinityModifier == pgtype.Infinity {
		return t.After(period.Lower.Time)
	}
	// If the period is finite, it contains everything between its bounds.
	return t.After(period.Lower.Time) && t.Before(period.Upper.Time)
}

/*============================================================================*/
/*=====*                              UUID                              *=====*/
/*============================================================================*/

// UnmarshalUUID : Implements the graphql.Unmarshaler interface for `UUID` GraphQL scalar and Go `pgtype.UUID` struct.
func UnmarshalUUID(value any) (pgtype.UUID, error) {
	uuid := pgtype.UUID{}
	if err := uuid.Set(value); err != nil {
		return uuid, errors.Errorf("`%v` is not a UUID", value)
	}
	return uuid, nil
}

func NewUUID() pgtype.UUID {
	return pgtype.UUID{Status: pgtype.Present, Bytes: uuid.New()}
}

func NullUUID() pgtype.UUID {
	return pgtype.UUID{Status: pgtype.Null}
}

func ParseUUID(v string) (pgtype.UUID, error) {
	data := pgtype.UUID{Status: pgtype.Present}
	return data, data.Set(v)
}

func FormatUUID(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

var (
	escaper   = strings.NewReplacer("9", "99", "-", "90", "_", "91")
	unescaper = strings.NewReplacer("99", "9", "90", "-", "91", "_")
)

func EncodeShortUUID(uuid pgtype.UUID) string {
	return escaper.Replace(base64.RawURLEncoding.EncodeToString(uuid.Bytes[:]))
}

func DecodeShortUUID(uuid string) (pgtype.UUID, error) {
	dec, err := base64.RawURLEncoding.DecodeString(unescaper.Replace(uuid))
	if err != nil {
		return pgtype.UUID{}, err
	}
	id := pgtype.UUID{}
	if err := id.Set(dec); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}
