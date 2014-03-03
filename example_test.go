package sh

import "testing"

func TestExample2(t *testing.T) {
	sh := NewSession()
	sh.Call("echo", []string{"hello", "example"})
	sh.Call("echo")
}
