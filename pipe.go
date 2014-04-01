package sh

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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

// unmarshal command output into xml
func (s *Session) UnmarshalXML(data interface{}) (err error) {
	bufrw := bytes.NewBuffer(nil)
	s.Stdout = bufrw
	if err = s.Run(); err != nil {
		return
	}
	return xml.NewDecoder(bufrw).Decode(data)
}

// start command
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
// only catch the last command error
func (s *Session) Wait() (err error) {
	for _, cmd := range s.cmds {
		err = cmd.Wait()
		wr, ok := cmd.Stdout.(*io.PipeWriter)
		if ok {
			wr.Close()
		}
	}
	return err
}

func (s *Session) Run() (err error) {
	if err = s.Start(); err != nil {
		return
	}
	return s.Wait()
}

func (s *Session) Output() (out []byte, err error) {
	oldout := s.Stdout
	defer func() {
		s.Stdout = oldout
	}()
	stdout := bytes.NewBuffer(nil)
	s.Stdout = stdout
	err = s.Run()
	out = stdout.Bytes()
	return
}
