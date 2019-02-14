package a

import "golang.org/x/xerrors"

func main() {
	var err error
	xerrors.Errorf("message: %w", err) // OK
	xerrors.Errorf("message:%w", err)  // want "invalid arguments"
	xerrors.Errorf(":%w", err)         // want "invalid arguments"
	xerrors.Errorf("%w", err)          // want "invalid arguments"
	xerrors.Errorf("%w", nil)          // want "invalid arguments"
	xerrors.Errorf("message: %w", nil) // want "invalid arguments"
	xerrors.Errorf("message: %w")      // want "invalid arguments"
}
