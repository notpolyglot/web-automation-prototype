package modules

import (
	"github.com/yuin/gopher-lua"
)

type HTTPClient struct {
}

var httpClientMethods = map[string]lua.LGFunction{
	"request": httpRequest,
}

func RegisterHTTPClientType(L *lua.LState) {
	mt := L.NewTypeMetatable()
	L.SetGlobal("HTTPClient", mt)
	L.SetField(mt, "new", L.NewFunction(newPerson))

}

func newHTTPClient(L *lua.LState) int {
	client := &HTTPClient{}

	ud := L.NewUserData()
	ud.Value = client

	return 1

}

func httpRequest(L *lua.LState) int {

}
