package cerrors

import (
	"errors"
)

/*============================================================================*/
/*=====*                           Base Error                           *=====*/
/*============================================================================*/

type iError interface {
	JSON() []byte
	CLI() string
}

type Error struct {
	Value iError
	next  *Error
}

func IsError(err error) *Error {
	var cerr *Error
	if errors.As(err, &cerr) {
		return cerr
	}
	return nil
}

func (e *Error) Error() string {
	return e.CLI()
}

func (e *Error) Append(err *Error) *Error {
	if err == nil {
		return e
	}
	err.next = e
	return err
}

func (e *Error) JSON() []byte {
	result := []byte{'['}
	for elm := e; elm != nil; elm = elm.next {
		result = append(result, elm.Value.JSON()...)
		if elm.next != nil {
			result = append(result, ',')
		}
	}
	result = append(result, ']')
	return result
}

func (e *Error) CLI() string {
	count := 1
	result := ""
	for elm := e; elm != nil; elm = elm.next {
		result = result + elm.Value.CLI()
		count = count + 1
		if elm.next != nil {
			result = result + "\n"
		}
	}
	return result
}
