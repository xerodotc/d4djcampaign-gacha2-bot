package session

import "errors"

var (
	ErrNotOK             = errors.New("not ok status")
	ErrNeverRolled       = errors.New("never rolled")
	ErrAlreadyGotSerial  = errors.New("already got serial")
	ErrRollLimitExceeded = errors.New("roll limit exceeded")
	ErrUnknown           = errors.New("unknown error")
)
