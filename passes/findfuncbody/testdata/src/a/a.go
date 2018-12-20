package a

func main() { // want "body is a.go:3"
	func() { // want "body is a.go:4"
	}()

	f := func() { // want "body is a.go:7"
	}

	f()        // want "body is a.go:7"
	g()        // want "body is decls.go:3"
	h()        // want "body is decls.go:6"
	new(T).M() // want "body is decls.go:13"
	T{}.F()    // want "body is decls.go:14"
	new(T).F() // want "body is decls.go:14"

	f = func() { // want "body is a.go:17"
	}
	f() // want "body is a.go:17"
}
