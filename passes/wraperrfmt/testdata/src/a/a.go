package a

import "golang.org/x/xerrors"

func main() {
	var err error
	xerrors.Errorf("message: %w", err)                   // OK
	xerrors.Errorf("message:%w", err)                    // want "invalid arguments"
	xerrors.Errorf(":%w", err)                           // want "invalid arguments"
	xerrors.Errorf("%w", err)                            // want "invalid arguments"
	xerrors.Errorf("%w", nil)                            // want "invalid arguments"
	xerrors.Errorf("message: %w")                        // want "invalid arguments"
	xerrors.Errorf("message: %w", nil)                   // want "invalid arguments"
	xerrors.Errorf("message: %w", nil, err)              // OK
	xerrors.Errorf("message: %w", []interface{}{err}...) // Unsupport
	args := []interface{}{err}
	xerrors.Errorf("message: %w", args...) // Unsupport
}
