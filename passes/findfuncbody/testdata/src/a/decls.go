package a

func g() {
}

func h() func() {
	return func() { // want "body is decls.go:7"
	}
}

type T struct{}

func (_ *T) M() {
}
