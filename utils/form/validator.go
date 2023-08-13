package form

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"time"

	cerrors "movies/utils/cerrors"

	en "github.com/go-playground/locales/en"
	fr "github.com/go-playground/locales/fr"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_t "github.com/go-playground/validator/v10/translations/en"
	fr_t "github.com/go-playground/validator/v10/translations/fr"
	pgtype "github.com/jackc/pgtype"
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	language "golang.org/x/text/language"
)

/*===========================================================================*/
/*=====*                              Init                             *=====*/
/*===========================================================================*/

var (
	validate *validator.Validate = initValidator()
)

func translator(lang language.Tag) ut.Translator {
	uni := ut.New(en.New(), en.New(), fr.New())
	trans, _ := uni.GetTranslator(lang.String())
	return trans
}

func initValidator() *validator.Validate {
	v := validator.New()

	if err := en_t.RegisterDefaultTranslations(v, translator(language.English)); err != nil {
		log.Fatal(err)
	}
	if err := fr_t.RegisterDefaultTranslations(v, translator(language.French)); err != nil {
		log.Fatal(err)
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register pgtype
	v.RegisterCustomTypeFunc(
		func(field reflect.Value) any {
			if valuer, ok := field.Interface().(driver.Valuer); ok {
				val, err := valuer.Value()
				if err == nil {
					return val
				}
			}
			return nil
		},
		pgtype.Bool{}, pgtype.Date{}, pgtype.Daterange{},
		pgtype.JSON{}, pgtype.JSONB{},
		pgtype.UUID{},
		pgtype.Timestamp{}, pgtype.Timestamptz{},
		pgtype.Tstzrange{}, pgtype.TextArray{},
		pgtype.Point{}, numeric.Numeric{},
	)

	// Custom validation
	Alphanumdot(v)
	Hexanumdot(v)
	Enum(v)

	return v
}

/*===========================================================================*/
/*=====*                              API                              *=====*/
/*===========================================================================*/

func GetValidator() *validator.Validate {
	return validate
}

func ValidateStruct(obj any, fieldPath bool) error {
	lang := language.English
	validator := GetValidator()
	err := validator.Struct(obj)
	return formatErrors(lang, err, fieldPath)
}

func ValidateJSON(r io.Reader, obj any) error {
	if err := json.NewDecoder(r).Decode(obj); err != nil {
		return errorJSON(err)
	}
	if e := ValidateStruct(obj, false); e != nil {
		return e
	}
	return nil
}

/*===========================================================================*/
/*=====*                            Internal                           *=====*/
/*===========================================================================*/

func errorJSON(err error) error {
	if errors.Is(err, io.EOF) {
		return cerrors.NewString("JSON Body not found")
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return cerrors.NewString("JSON Syntax error")
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return cerrors.NewString("JSON Field type error")
	}

	var unsupportedTypeError *json.UnsupportedTypeError
	if errors.As(err, &unsupportedTypeError) {
		return cerrors.NewString("JSON Unsupported type error")
	}

	var unsupportedValueError *json.UnsupportedValueError
	if errors.As(err, &unsupportedValueError) {
		return cerrors.NewString("JSON Unsupported value error")
	}

	var invalidUnmarshalError *json.InvalidUnmarshalError
	if errors.As(err, &invalidUnmarshalError) {
		return cerrors.NewString("JSON Invalide nnmarshal error")
	}

	var parseError *time.ParseError
	if errors.As(err, &parseError) {
		return cerrors.NewString(fmt.Sprintf("JSON Cannot parse `%s`", strings.Trim(parseError.Value, "\"")))
	}

	panic(err)
}

func formatErrors(lang language.Tag, errs error, fieldPath bool) error {
	if errs == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(errs, &validationErrors) {
		panic("Validator error should not be nil. 1")
	}

	var er *cerrors.Error
	for _, err := range validationErrors {
		var field string
		if fieldPath {
			field = err.Namespace()
		} else {
			field = err.Field()
		}
		e := cerrors.NewValidation(err.Tag(), field, err.Translate(translator(lang)), err.Value())
		er = er.Append(e)
	}
	if er == nil {
		panic("Validator error should not be nil. 2")
	}
	return er
}
