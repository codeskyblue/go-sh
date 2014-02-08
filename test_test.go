package sh

import "testing"

var s = NewSession()

type T struct{ *testing.T }

func NewT(t *testing.T) *T { return &T{t} }

func (t *T) checkTest(exp string, arg string, result bool) {
	r := s.Test(exp, arg)
	if r != result {
		t.Errorf("test -%s %s, %s != %s", exp, arg, r, result)
	}
}
func TestTest(i *testing.T) {
	t := NewT(i)
	t.checkTest("d", "../go-sh", true)
	t.checkTest("d", "../yymm", false)

	t.checkTest("f", "./sh.go", true)
	t.checkTest("f", "../go-sh", false)
	t.checkTest("f", "./yymm", false)

	//	t.checkTest("x", "lksjdf", false)
}
