package errortypes

import (
	"github.com/pritunl/tools/errors"
)

type ParseError struct {
	errors.StackError
}

type RequestError struct {
	errors.StackError
}
