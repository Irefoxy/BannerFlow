package errors

import (
	"errors"
	"fmt"
)

var (
	ErrorBadRequest           = errors.New("request is invalid")
	ErrorAuthenticationFailed = errors.New("authentication failed")
	ErrorNoPermission         = errors.New("error no permission")

	ErrorInRequestBody = fmt.Errorf("%w: error in request body", ErrorBadRequest)
	ErrorInParam       = fmt.Errorf("%w: error in param", ErrorBadRequest)
	ErrorNoToken       = fmt.Errorf("%w: error no token", ErrorAuthenticationFailed)
)
