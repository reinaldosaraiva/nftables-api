package errors

import "errors"

// ErrRecordNotFound é um erro personalizado para quando um registro não é encontrado.
var ErrRecordNotFound = errors.New("record not found")
