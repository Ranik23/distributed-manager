package errors

import "fmt"


var (
	ErrOpenFailure = fmt.Errorf("failed to open config file")
	ErrDecodeFailure = fmt.Errorf("failed to decode config file")
)