package form

import (
	"log"
	"reflect"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	language "golang.org/x/text/language"
)

/*============================================================================*/
/*=====*                        Custom Validator                        *=====*/
/*============================================================================*/

type translation struct {
	language    language.Tag
	translation string
}

type customValidator struct {
	Name        string
	Validate    func(fl validator.FieldLevel) bool
	Translation []translation
}

func (v customValidator) Register(vld *validator.Validate) {
	if err := vld.RegisterValidation(v.Name, v.Validate); err != nil {
		log.Fatal(err)
	}

	for _, t := range v.Translation {
		err := vld.RegisterTranslation(
			v.Name,
			translator(t.language),
			func(ut ut.Translator) (err error) {
				if err = ut.Add(v.Name, t.translation, true); err != nil {
					return
				}
				return
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), reflect.ValueOf(fe.Value()).String(), fe.Param())
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

/*============================================================================*/
/*=====*                           Alphanumdot                          *=====*/
/*============================================================================*/

func Alphanumdot(v *validator.Validate) {
	reg := regexp.MustCompile(`^[a-zA-Z0-9.]+$`)

	cv := customValidator{
		Name: "alphanumdot",
		Validate: func(fl validator.FieldLevel) bool {
			return reg.MatchString(fl.Field().String())
		},
		Translation: []translation{{
			language:    language.English,
			translation: "`{0}` is not valid, must be alphanumeric & dot characters only",
		}},
	}

	cv.Register(v)
}

func Hexanumdot(v *validator.Validate) {
	reg := regexp.MustCompile(`^[a-fA-F0-9.]+$`)

	cv := customValidator{
		Name: "hexanumdot",
		Validate: func(fl validator.FieldLevel) bool {
			return reg.MatchString(fl.Field().String())
		},
		Translation: []translation{{
			language:    language.English,
			translation: "`{0}` is not valid, must be hexanumeric & dot characters only",
		}},
	}

	cv.Register(v)
}

func Enum(v *validator.Validate) {
	cv := customValidator{
		Name: "enum",
		Validate: func(fl validator.FieldLevel) bool {
			type isValid interface {
				IsValid() bool
			}
			if valuer, ok := fl.Field().Interface().(isValid); ok {
				return valuer.IsValid()
			}
			return false
		},
		Translation: []translation{{
			language:    language.English,
			translation: "`{0}` is not a valid enum",
		}},
	}

	cv.Register(v)
}
