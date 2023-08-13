package form

import (
	"context"
	"unicode"

	errors "emperror.dev/errors"
	modifiers "github.com/go-playground/mold/v4/modifiers"
	pgtype "github.com/jackc/pgtype"
	lo "github.com/samber/lo"
	runes "golang.org/x/text/runes"
	transform "golang.org/x/text/transform"
	norm "golang.org/x/text/unicode/norm"
)

var conform = modifiers.New()

func ConformStruct(obj any) error {
	return conform.Struct(context.Background(), obj)
}

// Conform pgtype.Status
//
// Pgtype is flawed, when it's undefined, insertion fail.
// Until we migrate to pgx v5, we need to conform it.
func RemoveUndefined(obj *pgtype.Status) {
	if obj == nil {
		return
	} else if *obj == pgtype.Undefined {
		*obj = pgtype.Null
	}
}

// ShortenString: Shorten a string to a given length
func ShortenString(str string, length int) string {
	return string(lo.Slice([]rune(str), 0, length))
}

// RemoveDiacritics: Remove diacritics from a string
func RemoveDiacritics(str string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, str)
	if err != nil {
		return result, errors.WithStack(err)
	}
	return result, nil
}
