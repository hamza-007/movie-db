package cerrors

import (
	"encoding/json"
	"fmt"
	"log"

	pgtype "github.com/jackc/pgtype"
)

type stringError struct {
	Message string `json:"message"`
}

type scanner interface {
	EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error)
}

func NewString(str string, args ...any) *Error {
	arr := []any{}
	for _, arg := range args {
		switch f := arg.(type) {
		case scanner:
			buf, err := f.EncodeText(nil, nil)
			if err != nil {
				log.Fatal(err)
			}
			arr = append(arr, buf)
		default:
			arr = append(arr, f)
		}
	}
	return &Error{
		Value: stringError{Message: fmt.Sprintf(str, arr...)},
		next:  nil,
	}
}

func (e stringError) JSON() []byte {
	json, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return json
}

func (e stringError) CLI() string {
	return e.Message
}
