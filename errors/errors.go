package errors

import (
	"errors"
	"github.com/ngamux/ngamux/constants"
)

type NgamuxErr error

var (
	NotFound         NgamuxErr = errors.New(constants.NotFoundMessage)
	MethodNotAllowed NgamuxErr = errors.New(constants.MethodNotAllowedMessage)
)
