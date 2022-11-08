# goagi: Golang library to build agi/fastagi applications

[![Build Status](https://github.com/staskobzar/goagi/actions/workflows/ci.yml/badge.svg)](https://github.com/staskobzar/goagi/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/staskobzar/goagi/branch/master/graph/badge.svg)](https://codecov.io/gh/staskobzar/goagi)
[![CodeFactor](https://www.codefactor.io/repository/github/staskobzar/goagi/badge)](https://www.codefactor.io/repository/github/staskobzar/goagi)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/staskobzar/goagi/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/staskobzar/goagi/?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/staskobzar/goagi)](https://goreportcard.com/report/github.com/staskobzar/goagi)
![GitHub](https://img.shields.io/github/license/staskobzar/goagi)



Simple library that helps to build AGI scripts or FastAGI servers with Go.
```go
import "github.com/staskobzar/goagi"
```

API documentation [link is here](docs/api.md).

## Usage FastAGI

AGI object is created with ```New``` [method](docs/api.md#func-new) with three arguments:
- Reader [interface](docs/api.md#type-reader)
- Writer [interface](docs/api.md#type-response)
- Debugger [interface](docs/api.md#type-debugger)

Debugger interface is required only for debugging and usually ```nil```. See below for more details.

Reader and Writer are interfaces are any objects that implement ```Read```/```Write``` methods
and can be ```net.Conn```, ```tls.Conn```, ```os.Stdin```, ```os.Stdout``` or any other,
for example from packages ```strings```, ```bufio```, ```bytes```.

```New``` method will read AGI session setup environment variables and provides interface
to AGI commands. AGI environment variables and arguments can be accessed with 
methods [```Env```](docs/api.md#func-agi-env) and [```EnvArgs```](docs/api.md#func-agi-envargs).
If AGI channel receives HANGUP message, the session will be marked as hungup. Hangup status
can be checked by method [```IsHungup```](docs/api.md#func-agi-ishungup).

### Usage example for AGI:
```go
	import (
		"github.com/staskobzar/goagi"
		"os"
	)

	agi, err := goagi.New(os.Stdin, os.Stdout, nil)
	if err != nil {
		panic(err)
	}
	agi.Verbose("Hello World!")
```

### Fast AGI example:
```go
	ln, err := net.Listen("tcp", "127.0.0.1:4573")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
			agi, err := goagi.New(conn, conn, nil)
			if err != nil {
				panic(err)
			}
			agi.Verbose("Hello World!")
		}(conn)
	}
```

See working examples in [examples/] folder.

Index of methods that implements AGI commands [see here.](docs/api.md)

## Commands Response interface

Every AGI command method return interface [```Response```](docs/api.md#type-response).
This is interface that provides access to AGI response values.
Success response example:
```
200 result=1 (speech) endpos=9834523 results=5
```
Fail response example:
```
511 Command Not Permitted on a dead channel or intercept routine
```
There are success  and error responses. 
Response interface implements following methods:

* ```Code() int```: response code: 200, 510 etc.
* ```RawResponse() string```: returns full text of AGI response.
* ```Result() int```: returns value of result= field.
* ```Value() string```: returns response value field that comes in parentheses. For exmpale "(timeout)" in response "200 result=1 (timeout)".
* ```Data() string```: returns text for error responses and dtmf values for command like GetData.
* ```EndPos() int64```: returns value for endpos= field.
* ```Digit() string```: return digit from digit= field.
* ```SResults() int```: return value for results= field.

## Debugger

Interface that provides debugging capabilities with configurable output.

Example of usage:

```go
	dbg := logger.New(os.Stdout, "myagi:", log.Lmicroseconds)
	r, w := net.Pipe()
	agi, err := goagi.New(r, w, dbg)
```
