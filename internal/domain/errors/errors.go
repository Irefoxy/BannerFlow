package e

import (
	"errors"
	"fmt"
)

var (
	ErrorInternal             = errors.New("internal service error")
	ErrorBadRequest           = errors.New("request is invalid")
	ErrorAuthenticationFailed = errors.New("authentication failed")
	ErrorNoPermission         = errors.New("error no permission")
	ErrorNotFound             = errors.New("banner not found")

	ErrorFailedToConnect = fmt.Errorf("%w: failed to connect", ErrorInternal)
	ErrorConflict        = fmt.Errorf("%w: banner already exists", ErrorBadRequest)
	ErrorInRequestBody   = fmt.Errorf("%w: error in request body", ErrorBadRequest)
	ErrorInParam         = fmt.Errorf("%w: error in param", ErrorBadRequest)
	ErrorNoToken         = fmt.Errorf("%w: error no token", ErrorAuthenticationFailed)
)
