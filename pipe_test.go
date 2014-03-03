package sh

import "testing"

func TestUnmarshalJSON(t *testing.T) {
	var a int
	s := NewSession()
	s.ShowCMD = true
	err := s.Command("echo", []string{"1"}).UnmarshalJSON(&a)
	if err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect a tobe 1, but got %d", a)
	}
}
