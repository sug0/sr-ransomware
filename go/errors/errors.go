package errors

import (
    "fmt"
    goerrors "errors"
)

func Is(err, target error) bool {
    return goerrors.Is(err, target)
}

func Wrap(pkg, reason string, err error) error {
    if err != nil {
        return fmt.Errorf("%s: %s: %w", pkg, reason, err)
    }
    return nil
}

func Unwrap(err error) error {
    return goerrors.Unwrap(err)
}
