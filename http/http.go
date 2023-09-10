package http

import (
	"github.com/yuin/gopher-lua"
)

// i dont know how to handle proxies.
// the proxy function is set in the transport, and the docs recommend not to create transport structs on demand
// maybe its fine to make 1 transport per thread however i need to experiment
// although the structs have caching, which may cause problems with bot detection etc, so i think it would be bette to have 1 client & transport per script/job
func Loader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"myfunc": myfunc,
}

func myfunc(L *lua.LState) int {
	return 0
}
