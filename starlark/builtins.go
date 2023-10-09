package starlark

import "go.starlark.net/starlark"

var StarlarkBuiltIns = starlark.StringDict{
	"httpClient": newHTTPClient,
}
