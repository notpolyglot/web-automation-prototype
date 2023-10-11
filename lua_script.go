package main

import (
	// "log"
	// "go.starlark.net/starlark"
	"bufio"
	"errors"
	// "errors"
	"fmt"
	"os"
	"strings"
	"sync"

	libs "github.com/vadv/gopher-lua-libs"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
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

// to cache the exec func
type stateAndFunc struct {
	s *lua.LState
	f lua.LValue
}

type lStatePool struct {
	m     sync.Mutex
	saved []*stateAndFunc

	bytecode *lua.FunctionProto
}

func (pl *lStatePool) Get() *stateAndFunc {
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

func (pl *lStatePool) New() *stateAndFunc {
	L := lua.NewState()

	libs.Preload(L)

	DoCompiledFile(L, pl.bytecode)

	F := L.GetGlobal("execute")
	return &stateAndFunc{
		s: L,
		f: F,
	}
}

func (pl *lStatePool) Put(L *stateAndFunc) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *lStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.s.Close()
	}
}

type LuaScript struct {
	pool       *WorkerPool
	lPool      *lStatePool
	scriptPath string
}

type luaResponse struct {
	Capture  []any
	Counters []string
	Success  bool
	Retry    bool
}

func (s *LuaScript) execute(datum any) (r ExecResponse) {
	L := s.lPool.Get()
	defer s.lPool.Put(L)

	err := L.s.CallByParam(lua.P{
		Fn:      L.f,
		NRet:    1,
		Protect: true,
	}, lua.LString(datum.(string)))

	if err != nil {
		r.Error = err
		return
	}

	ret := L.s.Get(-1)
	L.s.Pop(1)

	//dont know if this will cause problems.
	//does allow the user to overwrite datum maybe, and error
	// var luaResp luaResponse
	if tbl, ok := ret.(*lua.LTable); ok {
		if err := gluamapper.Map(tbl, &r); err != nil {
			r.Error = err
			return
		}
	} else {
		r.Error = errors.New("execute function returned wrong type")
		return
	}
	return

}

func (s *LuaScript) print() {
	fmt.Print("\033[H\033[2J")

	for k, v := range s.pool.counters {
		fmt.Printf("%s: %d\n", strings.Title(k), v)
	}

	fmt.Printf("RPM: %v\n", s.pool.RPM())

}

func NewLuaScript(path string, size int, opts ...OptFunc) *LuaScript {
	s := &LuaScript{}
	s.scriptPath = path
	s.pool = NewWorkerPool(s.execute, size, append(opts, WithPrintFunc(s.print))...)

	s.lPool = &lStatePool{
		saved: make([]*stateAndFunc, 0, size),
	}

	return s
}

func (s *LuaScript) Start() error {

	bc, err := CompileLua(s.scriptPath)
	if err != nil {
		return err
	}

	s.lPool.bytecode = bc

	s.pool.Start()
	return nil
}
