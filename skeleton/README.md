# skeleton 

skeleton is create skeleton codes for golang.org/x/tools/go/analysis.

## Insall

```
$ go get -u github.com/tenntenn/gosa/skeleton
```

## How to use

### Create a skeleton codes in GOPATH

```
$ skeleton pkgname
pkgname
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Create a skeleton codes with import path

```
$ skeleton -path="github.com/tenntenn/pkgname"
pkgname
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Create cmd directory

```
$ skeleton -cmd pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```
