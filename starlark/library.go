package starlark

import "go.starlark.net/starlark"

var httpMethods = map[string]*starlark.Builtin{
	"Client": starlark.NewBuiltin("Client", newHTTPClient)
}
