package sh

import (
	"os"
	"path/filepath"
	"reflect"
)

func filetest(name string, modemask os.FileMode) (match bool, err error) {
	fi, err := os.Stat(name)
	if err != nil {
		return
	}
	match = (fi.Mode() & modemask) == modemask
	return
}

func (s *Session) pwd() string {
	dir := s.inj.Get(reflect.TypeOf(Dir(""))).String()
	if dir == "" {
		dir, _ = os.Getwd()
	}
	return dir
}

func (s *Session) abspath(name string) string {
	return filepath.Join(s.pwd(), name)
}

func (s *Session) Test(expression string, argument string) bool {
	var err error
	var fi os.FileInfo
	fi, err = os.Stat(s.abspath(argument))
	switch expression {
	case "d":
		return err == nil && fi.IsDir()
	case "f":
		return err == nil && fi.Mode().IsRegular()
		// case "x":
		//	 return err == nil && fi.Mode()&os.ModeExclusive != 0
		//case "h", "L":
		//	return err == nil && fi.Mode()&os.ModeSymlink != 0
	}
	return false
}
