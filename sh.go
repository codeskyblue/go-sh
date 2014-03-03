package sh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/codegangsta/inject"
)

type Return struct {
	Stdout string
	Stderr string
}

func (r *Return) String() string {
	return r.Stdout
}

func (r *Return) Trim() string {
	return strings.TrimSpace(r.Stdout)
}

func Capture(a ...interface{}) (ret *Return, err error) {
	s := NewSession()
	return s.Capture(a...)
}

type Dir string

type Session struct {
	inj     inject.Injector
	alias   map[string][]string
	cmds    []*exec.Cmd
	started bool
	Env     map[string]string
	Stdout  io.Writer
	Stderr  io.Writer
	ShowCMD bool // enable for debug
}

func (s *Session) writePrompt(args ...interface{}) {
	var ps1 = fmt.Sprintf("[golang-sh]$")
	args = append([]interface{}{ps1}, args...)
	fmt.Fprintln(s.Stderr, args...)
}

func NewSession(a ...interface{}) *Session {
	env := make(map[string]string)
	for _, key := range []string{"PATH"} {
		env[key] = os.Getenv(key)
	}
	s := &Session{
		inj:    inject.New(),
		alias:  make(map[string][]string),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    env,
	}
	dir := Dir("")
	args := []string{}
	s.inj.Map(env).Map(dir).Map(args)
	for _, v := range a {
		s.inj.Map(v)
	}
	return s
}

func (s *Session) Alias(alias, cmd string, args ...string) {
	v := []string{cmd}
	v = append(v, args...)
	s.alias[alias] = v
}

func isInInj(inj inject.Injector, vType reflect.Type, v interface{}) bool {
	return inj.Get(vType).Kind() != reflect.Invalid && reflect.TypeOf(v).Kind() == vType.Kind()
}

func (s *Session) Command(a ...interface{}) *Session {
	var cmd string
	var first = true
	var args = make([]string, 0)
	var sType = reflect.TypeOf("")

	// extract arg0, args ...
	var sepIndex int = -1
	for i, v := range a {
		sepIndex = i
		if reflect.TypeOf(v) == sType {
			if first {
				first = false
				cmd = v.(string)
				continue
			}
			args = append(args, v.(string))
			continue
		}
		break
	}

	//fmt.Println("left=(", sepIndex, a, "args=", a[sepIndex:])
	//fmt.Println("cmd=", cmd)

	// extract sh.Dir ...
	s.inj.Map(args)
	for _, v := range a[sepIndex:] { // args can be replaced here
		s.inj.Map(v)
	}
	s.inj.Map(cmd)
	s.inj.Invoke(s.appendCmd)
	return s
}

func (s *Session) Call(a ...interface{}) error {
	return s.Command(a...).Run()
}

/*
func (s *Session) Exec(cmd string, args ...string) error {
	return s.Call(cmd, args)
}
*/

func (s *Session) Capture(a ...interface{}) (ret *Return, err error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	oldout, olderr := s.Stdout, s.Stderr
	s.Stdout, s.Stderr = stdout, stderr
	err = s.Call(a...)
	s.Stdout, s.Stderr = oldout, olderr

	ret = new(Return)
	ret.Stdout = string(stdout.Bytes())
	ret.Stderr = string(stderr.Bytes())
	return
}

func (s *Session) Set(a ...interface{}) *Session {
	for _, v := range a {
		s.inj.Map(v)
	}
	return s
}

func (s *Session) appendCmd(cmd string, args []string, cwd Dir) {
	if s.started {
		s.started = false
		s.cmds = make([]*exec.Cmd, 0)
	}
	envs := make([]string, 0, len(s.Env))
	for k, v := range s.Env {
		envs = append(envs, k+"="+v)
	}
	// add system default env
	for _, line := range os.Environ() {
		exists := false
		for k, _ := range s.Env {
			if strings.HasPrefix(line, k) {
				exists = true
				break
			}
		}
		if !exists {
			envs = append(envs, line)
		}
	}

	v, ok := s.alias[cmd]
	if ok {
		cmd = v[0]
		args = append(v[1:], args...)
	}
	c := exec.Command(cmd, args...)
	c.Env = envs
	c.Dir = string(cwd)
	s.cmds = append(s.cmds, c)
}
