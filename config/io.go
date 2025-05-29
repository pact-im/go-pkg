package config

import (
	"errors"
	"io/fs"
)

var skipErrors = []error{
	fs.ErrNotExist,
	fs.ErrInvalid,
	fs.ErrPermission,
}

func shouldSkipError(err error) bool {
	for _, errSkip := range skipErrors {
		if errors.Is(err, errSkip) {
			return true
		}
	}
	return false
}
