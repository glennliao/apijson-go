package tt

import (
	"fmt"
	"testing"
)

type A struct {
	Name string
}

func (a *A) M() string {
	return "A:" + a.Name
}

type B struct {
	A
}

func test(a A) {
	fmt.Print(a.M())
}

func TestName(t *testing.T) {
	b := &B{}

	b.Name = "asdsad"
	test(b.A)
}
