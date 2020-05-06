package errors

import (
    "fmt"
    goerrors "errors"
)

func New(pkg, msg string) error {
    return goerrors.New(pkg + ": " + msg)
}

func Is(err, target error) bool {
    return goerrors.Is(err, target)
}

func Wrap(pkg, reason string, err error) error {
    return fmt.Errorf("%s: %s: %w", pkg, reason, err)
}

func WrapIfNotNil(pkg, reason string, err error) error {
    if err != nil {
        return Wrap(pkg, reason, err)
    }
    return nil
}

func Unwrap(err error) error {
    return goerrors.Unwrap(err)
}
