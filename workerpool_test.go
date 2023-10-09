package main

import (
	// "log"
	// "go.starlark.net/starlark"
	"fmt"
	"strings"
	"testing"
	"time"
	// "time"
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

func TestWorkerPoolWithScriptInterface(t *testing.T) {
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

	if s.pool.counters["success"] != 5 {
		t.Error("wrong amount of success")
	}

	if s.pool.counters["fail"] != 5 {
		t.Error("wrong amount of fail")
	}

}
