package errTrace

import "fmt"

func Wrap(packageName string, funcName string, err *error, prevStack *error) error {
	return fmt.Errorf("[%s: %s]: %w\n%s", packageName, funcName, *err, *prevStack)
}
