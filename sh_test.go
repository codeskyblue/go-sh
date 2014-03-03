package sh

import (
	"runtime"
	"testing"
)

func TestAlias(t *testing.T) {
	session := NewSession()
	session.Alias("gr", "echo", "hi")
	ret, err := session.Capture("gr", []string{"sky"})
	if err != nil {
		t.Error(err)
	}
	if ret.Trim() != "hi sky" {
		t.Errorf("expect 'hi sky' but got:%s", ret)
	}
}

func TestCapture(t *testing.T) {
	r, err := Capture("echo", []string{"hello"})
	if err != nil {
		t.Error(err)
	}
	_ = r
	if r.Trim() != "hello" {
		t.Errorf("expect hello, but got %s", r.Trim())
	}
}

func TestSession(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("ignore test on windows")
		return
	}
	session := NewSession()
	session.Set(Dir("/"))
	session.ShowCMD = true
	err := session.Call("pwd")
	if err != nil {
		t.Error(err)
	}
	ret, err := session.Capture("pwd")
	if err != nil {
		t.Error(err)
	}
	if ret.Trim() != "/" {
		t.Errorf("expect /, but got %s", ret.Trim())
	}
}

/*
	#!/bin/bash -
	#
	export PATH=/usr/bin:/bin
	alias ll='ls -l'
	cd /usr
	if test -d "local"
	then
		ll local | awk '{print $1, $NF}' | grep bin
	fi
*/
func TestExample(t *testing.T) {
	s := NewSession()
	s.ShowCMD = true
	s.Env["PATH"] = "/usr/bin:/bin"
	s.Set(Dir("/usr"))
	s.Alias("ll", "ls", "-l")
	//s.Stdout = nil
	if s.Test("d", "local") {
		//s.Command("ll", []string{"local"}).Command("awk", []string{"{print $1, $NF}"}).Command("grep", []string{"bin"}).Run()
		s.Command("ll", "local").Command("awk", "{print $1, $NF}").Command("grep", "bin").Run()
	}
}
