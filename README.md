## go-sh
[![Build Status](https://drone.io/github.com/shxsun/go-sh/status.png)](https://drone.io/github.com/shxsun/go-sh/latest)
[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/shxsun/go-sh)

So what is go-sh. Sometimes we need to write some shell scripts, but shell scripts is not good at cross platform, but golang is good at that. Is there a good way to use golang to write scripts like shell? Use go-sh we can do it now.

go-sh support some shell futures.

* shell session
* `export`: env
* `alias`: like alias ll='ls -l'
* `cd`: remember current dir
* pipe
* `test`: this is shell build command, very usefull(support -d and -f)

Example is always important. I will show you how to use it.



First give you a full example, I will explain every command below.

	session := sh.NewSession()
	session.Env["PATH"] = "/usr/bin:/bin"
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Alias("ll", "ls", "-l")
	session.ShowCMD = true // enable for debug
	var err error
	err = session.Call("ll", "/")
	if err != nil {
		log.Fatal(err)
	}
	ret, err := session.Capture("pwd", sh.Dir("/home")) # wraper of session.Call
	if err != nil {
		log.Fatal(err)
	}
	# ret is "/home\n"
	fmt.Println(ret)

create a new Session

	session := sh.NewSession()

use alias like this

	session.Alias("ll", "ls", "-l") # like alias ll='ls -l'

set current env like this

	session.Env["BUILD_ID"] = "123" # like export BUILD_ID=123

set current directory

	session.Set(sh.Dir("/")) # like cd /

pipe is also supported

	session.Command("echo", "hello\tworld").Command("cut", "-f2")
	// output should be "world"
	session.Run()

test, the build in command support

	session.Test("d", "dir") // test dir
	session.Test("f", "file) // test regular file

with `Alias Env Set Call Capture Command` a shell scripts can be easily converted into golang program. below is a shell script.

	#!/bin/bash -
	#
	export PATH=/usr/bin:/bin
	alias ll='ls -l'
	cd /usr
	if test -d "local"
	then
		ll local | awk '{print $1, $NF}'
	fi

convert to golang, will be

	s := sh.NewSession()
	s.Env["PATH"] = "/usr/bin:/bin"
	s.Set(sh.Dir("/usr"))
	s.Alias("ll", "ls", "-l")
	if s.Test("d", "local") {
		s.Command("ll", "local").Command("awk", "{print $1, $NF}").Run()
	}

### contribute
If you love this project, star it which will encourage the coder. pull requests are welcomed, if you want to add some new fetures.

support the author: [alipay](https://me.alipay.com/goskyblue)

### thanks
this project is based on <http://github.com/codegangsta/inject>. thanks for the author.
