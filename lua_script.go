package main

import (
	// "log"
	// "go.starlark.net/starlark"
	"bufio"
	"fmt"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"os"
	"strings"
	"sync"
)

// CompileLua reads the passed lua file from disk and compiles it.
func CompileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func DoCompiledFile(L *lua.LState, proto *lua.FunctionProto) error {
	lfunc := L.NewFunctionFromProto(proto)
	L.Push(lfunc)
	return L.PCall(0, lua.MultRet, nil)
}

type lStatePool struct {
	m     sync.Mutex
	saved []*lua.LState
}

func (pl *lStatePool) Get() *lua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.New()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

func (pl *lStatePool) New() *lua.LState {
	L := lua.NewState()
	// setting the L up here.
	// load scripts, set global variables, share channels, etc...
	return L
}

func (pl *lStatePool) Put(L *lua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *lStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}

type LuaScript struct {
	pool       *WorkerPool
	lPool      *lStatePool
	bytecode   *lua.FunctionProto
	scriptPath string
}

func (s *LuaScript) execute(datum any) (r ExecResponse) {
	L := s.lPool.Get()
	defer s.lPool.Put(L)

	err := DoCompiledFile(L, s.bytecode)
	if err != nil {
		fmt.Println(err)
	}

	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal("execute"),
		NRet:    1,
		Protect: true,
	}, lua.LString(datum.(string)))

	if err != nil {
		r.Error = err
	}

	return
}

func (s *LuaScript) print() {
	// fmt.Print("\033[H\033[2J")

	for k, v := range s.pool.counters {
		fmt.Printf("%s: %d\n", strings.Title(k), v)
	}

	fmt.Printf("RPM: %v\n", s.pool.RPM())

}

func NewLuaScript(path string, size int, opts ...OptFunc) *LuaScript {
	s := &LuaScript{}
	s.scriptPath = path
	s.pool = NewWorkerPool(s.execute, size, WithPrintFunc(s.print))

	s.lPool = &lStatePool{
		saved: make([]*lua.LState, 0, size),
	}

	return s
}

func (s *LuaScript) Start() error {

	bc, err := CompileLua(s.scriptPath)
	if err != nil {
		return err
	}

	s.bytecode = bc

	s.pool.Start()
	return nil
}
