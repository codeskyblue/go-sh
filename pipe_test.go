package sh

import (
	"io"
	"os"
	"os/exec"
	"testing"
)

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

func TestPipe(t *testing.T) {
	s := NewSession()
	s.ShowCMD = true
	s.Call("echo", "hello")
	err := s.Command("echo", "hi").Command("cat", "-n").Start()
	if err != nil {
		t.Error(err)
	}
	err = s.Wait()
	if err != nil {
		t.Error(err)
	}
	out, err := s.Command("echo", []string{"hello"}).Output()
	if err != nil {
		t.Error(err)
	}
	if out != "hello\n" {
		t.Error("capture wrong output:", out)
	}
	s.Command("echo", []string{"hello\tworld"}).Command("cut", []string{"-f2"}).Run()
}

func TestPipeCommand(t *testing.T) {
	c1 := exec.Command("echo", "good")
	rd, wr := io.Pipe()
	c1.Stdout = wr
	c2 := exec.Command("cat", "-n")
	c2.Stdout = os.Stdout
	c2.Stdin = rd
	c1.Start()
	c2.Start()

	c1.Wait()
	wc, ok := c1.Stdout.(io.WriteCloser)
	if ok {
		wc.Close()
	}
	c2.Wait()
}
