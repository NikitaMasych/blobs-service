package config

import "github.com/pkg/errors"

func suppressStack(err error) error {
	if err == nil {
		return nil
	}
	return &suppressableErr{err, true}
}

type suppressableErr struct {
	err      error
	suppress bool
}

func (e suppressableErr) Error() string {
	return e.err.Error()
}

func (e suppressableErr) SuppressStack() bool {
	return e.suppress
}

func (e suppressableErr) Cause() error {
	return errors.Cause(e.err)
}
