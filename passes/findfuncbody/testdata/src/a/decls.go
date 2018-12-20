package a

func g() { // want "body is decls.go:3"
}

func h() func() { // want "body is decls.go:6"
	return func() { // want "body is decls.go:7"
	}
}

type T struct{ F func() }

func (t *T) M() { // want "body is decls.go:13"
	t.F = func() {} // want "body is decls.go:14"
}
