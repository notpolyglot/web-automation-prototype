package lualib

import (
	"github.com/yuin/gopher-lua"
)

type Response struct {
	Capture  []any
	Counters []string
	Success  bool
	Retry    bool
}

const luaResponseTypeName = "Response"

// Registers my person type to given L.
func RegisterResponseType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaResponseTypeName)
	L.SetGlobal("response", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newResponse))
	// methods
}

// Constructor
func newResponse(L *lua.LState) int {
	person := &Response{}
	ud := L.NewUserData()
	ud.Value = person
	L.SetMetatable(ud, L.GetTypeMetatable(luaResponseTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkResponse(L *lua.LState) *Response {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Response); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}

func responseSetRetry(L *lua.LState) int {
	p := checkResponse(L)
	p.Retry = true
	return 1
}

// Getter and setter for the Person#Name
// func personGetSetName(L *lua.LState) int {
// 	// p := checkResponse(L)
// 	// if L.GetTop() == 2 {
// 	// 	p.Name = L.CheckString(2)
// 	// 	return 0
// 	// }
// 	// L.Push(lua.LString(p.Name))
// 	return 1
// }
