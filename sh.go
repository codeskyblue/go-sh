package sh

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/codegangsta/inject"
)

/*
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

func Capture(name string, a ...interface{}) (ret *Return, err error) {
	s := NewSession()
	return s.Capture(name, a...)
}
*/

type Dir string

type Session struct {
	inj     inject.Injector
	alias   map[string][]string
	cmds    []*exec.Cmd
	dir     Dir
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

func NewSession() *Session {
	env := make(map[string]string)
	for _, key := range []string{"PATH"} {
		env[key] = os.Getenv(key)
	}
	s := &Session{
		inj:    inject.New(),
		alias:  make(map[string][]string),
		dir:    Dir(""),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    env,
	}
	return s
}

func Command(name string, a ...interface{}) *Session {
	s := NewSession()
	return s.Command(name, a...)
}

func (s *Session) Alias(alias, cmd string, args ...string) {
	v := []string{cmd}
	v = append(v, args...)
	s.alias[alias] = v
}

func (s *Session) Command(name string, a ...interface{}) *Session {
	var args = make([]string, 0)
	var sType = reflect.TypeOf("")

	// init cmd, args, dir, envs
	// if not init, program may panic
	s.inj.Map(name).Map(args).Map(s.dir).Map(map[string]string{})
	for _, v := range a {
		switch reflect.TypeOf(v) {
		case sType:
			args = append(args, v.(string))
		default:
			s.inj.Map(v)
		}
	}
	if len(args) != 0 {
		s.inj.Map(args)
	}
	s.inj.Invoke(s.appendCmd)
	return s
}

// combine Command and Run
func (s *Session) Call(name string, a ...interface{}) error {
	return s.Command(name, a...).Run()
}

/*
func (s *Session) Exec(cmd string, args ...string) error {
	return s.Call(cmd, args)
}
*/

/*
func (s *Session) Capture(name string, a ...interface{}) (ret *Return, err error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	oldout, olderr := s.Stdout, s.Stderr
	s.Stdout, s.Stderr = stdout, stderr
	err = s.Call(name, a...)
	s.Stdout, s.Stderr = oldout, olderr

	ret = new(Return)
	ret.Stdout = string(stdout.Bytes())
	ret.Stderr = string(stderr.Bytes())
	return
}
*/

func (s *Session) SetEnv(key, value string) *Session {
	s.Env[key] = value
	return s
}

func (s *Session) SetDir(dir string) *Session {
	s.dir = Dir(dir)
	return s
}

func newEnviron(env map[string]string, inherit bool) []string { //map[string]string {
	environ := make([]string, 0, len(env))
	if inherit {
		for _, line := range os.Environ() {
			for k, _ := range env {
				if strings.HasPrefix(line, k+"=") {
					goto CONTINUE
				}
			}
			environ = append(environ, line)
		CONTINUE:
		}
	}
	for k, v := range env {
		environ = append(environ, k+"="+v)
	}
	return environ
}

func (s *Session) appendCmd(cmd string, args []string, cwd Dir, env map[string]string) {
	if s.started {
		s.started = false
		s.cmds = make([]*exec.Cmd, 0)
	}
	for k, v := range s.Env {
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	environ := newEnviron(s.Env, true) // true: inherit sys-env
	v, ok := s.alias[cmd]
	if ok {
		cmd = v[0]
		args = append(v[1:], args...)
	}
	c := exec.Command(cmd, args...)
	c.Env = environ
	c.Dir = string(cwd)
	s.cmds = append(s.cmds, c)
}
