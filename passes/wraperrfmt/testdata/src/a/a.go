package a

import "golang.org/x/xerrors"

func main() {
	var err error
	xerrors.Errorf("message: %w", err)      // OK
	xerrors.Errorf("message:%w", err)       // want "unexpected format. format must end with ': %w'"
	xerrors.Errorf(":%w", err)              // want "unexpected format. format must end with ': %w'"
	xerrors.Errorf("%w", err)               // want "unexpected format. format must end with ': %w'"
	xerrors.Errorf("%w", nil)               // want "unexpected format. format must end with ': %w'"
	xerrors.Errorf("message: %w")           // want "unexpected format. format must end with ': %w'"
	xerrors.Errorf("message: %w", nil)      // want "unexpected format. format must end with ': %w'"
	xerrors.Errorf("message: %w", nil, err) // OK

	// Unsupport
	xerrors.Errorf("message: %w", []interface{}{err}...)
	args := []interface{}{err}
	xerrors.Errorf("message: %w", args...)
}
