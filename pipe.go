package sh

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
)

// unmarshal shell output to decode json
func (s *Session) UnmarshalJSON(data interface{}) (err error) {
	bufrw := bytes.NewBuffer(nil)
	s.Stdout = bufrw
	if err = s.Run(); err != nil {
		return
	}
	return json.NewDecoder(bufrw).Decode(data)
}

func (s *Session) Start() (err error) {
	s.started = true
	var rd *io.PipeReader
	var wr *io.PipeWriter
	var length = len(s.cmds)
	if s.ShowCMD {
		var cmds = make([]string, 0, 4)
		for _, cmd := range s.cmds {
			cmds = append(cmds, strings.Join(cmd.Args, " "))
		}
		s.writePrompt(strings.Join(cmds, " | "))
	}
	for index, cmd := range s.cmds {
		if index != 0 {
			cmd.Stdin = rd
		}
		if index != length {
			rd, wr = io.Pipe() // create pipe
			cmd.Stdout = wr
			cmd.Stderr = os.Stderr
		}
		if index == length-1 {
			cmd.Stdout = s.Stdout
			cmd.Stderr = s.Stderr
		}
		err = cmd.Start()
		if err != nil {
			return
		}
	}
	return
}

// Should be call after Start()
func (s *Session) Wait() (err error) {
	for _, cmd := range s.cmds {
		cmd.Wait()
		wr, ok := cmd.Stdout.(*io.PipeWriter)
		if ok {
			wr.Close()
		}
	}
	return
}

func (s *Session) Run() (err error) {
	if err = s.Start(); err != nil {
		return
	}
	return s.Wait()
}

func (s *Session) Output() (out string, err error) {
	oldout := s.Stdout
	defer func() {
		s.Stdout = oldout
	}()
	stdout := bytes.NewBuffer(nil)
	s.Stdout = stdout
	err = s.Run()
	out = string(stdout.Bytes())
	return
}
