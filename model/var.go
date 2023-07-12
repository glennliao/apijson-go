package model

type Var interface {
	String() string
	Int() int
	Scan(pointer interface{}, mapping ...map[string]string) error
}
