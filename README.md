# Unnaamed project
please help me name  it


# What is it

A workerpool/web automation toolkit which makes use of libraries such as gopher-lua and ants to make it easier to write fast tools.

## Write Lua

```lua
function execute()
  local http = require("http")
  local client = http.client()
  local request = http.request("GET", "https://google.com")
  local result, err = client:do_request(request)
  if err then error(err) end

  return {
    success = true,
    counters = {"fart"}
  } 
end
```

## Use as a library in Go
```go
package main

import (
	"fmt"
	"strings"
	"time"
)

type TestScript struct {
	pool *WorkerPool
}

type ExecData struct {
	email string
	pass  string
}

func DummyLogin(email string, pass string) bool {
	users := map[string]string{
		"pinkearwax@email.com": "real!!!123",
	}

	realpass, ok := users[email]
	if ok {
		if realpass == pass {
			return true
		}
	}
	return false
}

func (s TestScript) execute(datum any) (r ExecResponse) {
	time.Sleep(1 * time.Second)
	d := datum.(ExecData)

	r.Capture = []any{"example_value"}
	success := DummyLogin(d.email, d.pass)

	if success {
		r.Success = true
		return
	}

	return
}

func (s *TestScript) print() {
	fmt.Print("\033[H\033[2J")

	for k, v := range s.pool.counters {
		fmt.Printf("%s: %d\n", strings.Title(k), v)
	}

	fmt.Printf("RPM: %v\n", s.pool.RPM())

}

func NewTestScript() *TestScript {
	s := &TestScript{}

	s.pool = NewWorkerPool(s.execute, 1,
		WithPrintFunc(s.print),
		WithDatabase(NewFlatFileDatabase("test_log.txt", "line=%s thingamabob=%s")))
	s.pool.Start()

	return s
}

func main() {
	s := NewTestScript()

	for i := 0; i < 5; i++ {
		s.pool.DChan <- ExecData{
			email: "pinkearwax@email.com",
			pass:  "real!!!123",
		}
	}

	for i := 0; i < 5; i++ {
		s.pool.DChan <- ExecData{
			email: "pinkearwax@email.com",
			pass:  "wrongpassword",
		}
	}

	s.pool.Stop()
}
```

## Example Use cases
 - Fuzzing
 - Credential stuffing
 - Web scraping

## Features planned to be implemented (unordered)

 - Custom Lua libraries (currently using a package)
 - Libraries for web scraping
 - Debugger with tools such as a HTML renderer 
 - GUI made in Wails
 - PostgreSQL and SQLite3 database interfaces
 
