package internal

import "errors"

type Warning string

func (warning Warning) Error() string { return string(warning) }

func (Warning) Is(err error) bool {
	_, ok := err.(Warning)
	return ok
}

func IsWarning(err error) bool {
	return errors.Is(err, Warning(""))
}
