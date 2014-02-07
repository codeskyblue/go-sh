## go-sh
[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/shxsun/go-sh)

So what is go-sh. Sometimes I need to write some shell scripts, but shell scripts is not good at cross platform, but golang is good at that. Is there a good way to use golang to write scripts like shell? Use go-sh we can do it now.

go-sh support some shell futures.

* Session
* `export`: env
* `alias`: like alias ll='ls -l'
* `cd`: remember current dir
* pipe

Example is very important. I will show you how to use it.

run `echo hi world` in dir(/)

	sh.Capture("echo", []string{"hi", "world"}, sh.Dir("/"))

create a new Session

	session := sh.NewSession()
	session.Env["PATH"] = "/usr/bin:/bin"
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Alias("ll", "ls", "-l")
	var err error
	err = session.Call("ll", []string{"/"})
	if err != nil {
		log.Fatal(err)
	}
	ret, err := session.Capture("pwd", sh.Dir("/home"))
	if err != nil {
		log.Fatal(err)
	}
	# ret is "/home\n"
	fmt.Println(ret)

use alias like this

	session.Alias("ll", "ls", "-l") # like alias ll='ls -l'

set current env like this

	session.Env["BUILD_ID"] = "123" # like export BUILD_ID=123

set current directory

	session.Set(sh.Dir("/")) # like cd /

empty args filled in Call will call last command

	session.Exec("echo", []string{"hi"})
	session.Call() # will call echo hi again

pipe is also supported

	session.Command("echo", []string{"hello\tworld"}).Command("cut", []string{"-f2"})
	// output should be "world"
	session.Run()

with `Alias Env Set Call Capture Command` a shell scripts can be easily converted into golang program. below is a shell script.

	#!/bin/bash -
	#
	export PATH=/usr/bin:/bin
	alias ll='ls -l'
	ll | awk '{print $1}' | grep "^-rw"

convert to golang, will be

	s := sh.NewSession()
	s.Env["PATH"] = "/usr/bin:/bin"
	s.Alias("ll", "ls", "-l")
	s.Command("ll").Command("awk", []string{"'{print $1}'"}).Command("grep", []string{"^-rw"}).Run()

### contribute
If you love this project, star it which will encourage the coder. pull requests are welcomed, if you want to add some new fetures.

### thanks
this project is based on <http://github.com/codegangsta/inject>. thanks for the author.
