[![CI](https://github.com/gomodules/go-sh/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/gomodules/go-sh/actions/workflows/ci.yml)
[![PkgGoDev](https://pkg.go.dev/badge/gomodules.xyz/go-sh)](https://pkg.go.dev/gomodules.xyz/go-sh)

# go-sh

install: `go get gomodules.xyz/go-sh`

Pipe Example:

	package main

	import "gomodules.xyz/go-sh"

	func main() {
		sh.Command("echo", "hello\tworld").Command("cut", "-f2").Run()
	}

Because I like os/exec, `go-sh` is very much modelled after it. However, `go-sh` provides a better experience.

These are some of its features:

* keep the variable environment (e.g. export)
* alias support (e.g. alias in shell)
* remember current dir
* pipe command
* shell build-in commands echo & test
* timeout support
* run multiple concurrent leaf commands using a single pipe input

Examples are important:

	sh: echo hello
	go: sh.Command("echo", "hello").Run()

	sh: export BUILD_ID=123
	go: s = sh.NewSession().SetEnv("BUILD_ID", "123")

	sh: alias ll='ls -l'
	go: s = sh.NewSession().Alias('ll', 'ls', '-l')

	sh: (cd /; pwd)
	go: sh.Command("pwd", sh.Dir("/")).Run()

	sh: test -d data || mkdir data
	go: if ! sh.Test("dir", "data") { sh.Command("mkdir", "data").Run() }

	sh: cat first second | awk '{print $1}'
	go: sh.Command("cat", "first", "second").Command("awk", "{print $1}").Run()

	sh: count=$(echo "one two three" | wc -w)
	go: count, err := sh.Echo("one two three").Command("wc", "-w").Output()

	sh(in ubuntu): timeout 1s sleep 3
	go: c := sh.Command("sleep", "3"); c.Start(); c.WaitTimeout(time.Second) # default SIGKILL
	go: out, err := sh.Command("sleep", "3").SetTimeout(time.Second).Output() # set session timeout and get output)

	sh: echo hello | cat
	go: out, err := sh.Command("cat").SetInput("hello").Output()

	sh: cat # read from stdin
	go: out, err := sh.Command("cat").SetStdin(os.Stdin).Output()

	sh: ls -l > /tmp/listing.txt # write stdout to file
	go: err := sh.Command("ls", "-l").WriteStdout("/tmp/listing.txt")

If you need to keep env and dir, it is better to create a session

	session := sh.NewSession()
	session.SetEnv("BUILD_ID", "123")
	session.SetDir("/")
	# then call cmd
	session.Command("echo", "hello").Run()
	# set ShowCMD to true for easily debug
	session.ShowCMD = true

By default, pipeline returns error only if the last command exit with a non-zero status. However, you can also enable `pipefail` option like `bash`. In that case, pipeline returns error if any of the commands fail and for multiple failed commands, it returns the error of rightmost failed command.

	session := sh.NewSession()
	session.PipeFail = true
	session.Command("cat", "unknown-file").Command("echo").Run()

By default, pipelines's std-error is set to last command's std-error. However, you can also combine std-errors of all commands into pipeline's std-error using `session.PipeStdErrors = true`.

By default, pipeline returns error only if the last command exit with a non-zero status. However, you can also enable `pipefail` option like `bash`. In that case, pipeline returns error if any of the commands fail and for multiple failed commands, it returns the error of rightmost failed command.

	session := sh.NewSession()
	session.PipeFail = true
	session.Command("cat", "unknown-file").Command("echo").Run()


Designing a Command Chain to Run Multiple Concurrent Leaf Commands Using a Single Pipe Input. 

Features be like:
* **Input Sharing**: All leaf commands take the same input from a pipe.
* **Separate Environments**: Each leaf command runs with its own environment variables.
* **Output Aggregation**: Outputs from all commands are combined into a single result.
* **Error Handling**: Errors are collected and included in the output (e.g., shell or variable).
* **Timeouts**: Each command has same timeout and will apply simultaneously.

Below is an example of multiple concurrent leaf commands using a single pipe input

	s := sh.NewSession()
	s.ShowCMD = true
	s.Command("echo", "hello world").LeafCommand("xargs").LeafCommand("xargs")
	s.Run()

Below is an example of each leaf command runs with its own environment variables
    
    s := sh.NewSession()
	s.ShowCMD = true
	var args1,args2 []interface{}
	
	mp := make(map[string]string)
	mp["COMPANY_NAME"] = "APPSCODE"
	args1 = append(args1, "COMPANY_NAME")
	args1 = append(args1, mp)
	s.LeafCommand("printenv", args1...)

	mp["COMPANY_NAME"] = "GOOGLE"
	args2 = append(args2, "COMPANY_NAME")
	args2 = append(args2, mp)
	s.LeafCommand("printenv", args2...)
	
	s.Run()

for more information, it better to see docs.
[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/gomodules.xyz/go-sh)

### contribute
If you love this project, starring it will encourage the coder. Pull requests are welcome.

support the author: [alipay](https://me.alipay.com/goskyblue)

### thanks
this project is based on <http://github.com/codegangsta/inject>. thanks for the author.

# the reason to use Go shell
Sometimes we need to write shell scripts, but shell scripts are not good at working cross platform,  Go, on the other hand, is good at that. Is there a good way to use Go to write shell like scripts? Using go-sh we can do this now.
