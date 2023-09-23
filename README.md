<div align="center">
<h1>✨ GJS ✨</h1>
<p><strong>Comprehensive API for Webassembly environment</strong></p>

<p>
    <a href="https://github.com/chumaumenze/gjs/fork"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat" alt="Create a fork"></a>
    <a href="https://github.com/chumaumenze/gjs/actions"><img src="https://github.com/chumaumenze/gjs/actions/workflows/test.yml/badge.svg" alt="Github Actions"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg" alt="made-with-Go"></a>
    <a href="https://goreportcard.com/report/github.com/chumaumenze/gjs"><img src="https://goreportcard.com/badge/github.com/chumaumenze/gjs" alt="GoReportCard"></a>
    <a href="https://github.com/chumaumenze/gjs"><img src="https://img.shields.io/github/go-mod/go-version/chumaumenze/gjs.svg" alt="Go.mod version"></a>
    <a href="https://github.com/chumaumenze/gjs/blob/master/LICENCE"><img src="https://img.shields.io/github/license/chumaumenze/gjs.svg" alt="LICENCE"></a>
    <a href="https://github.com/chumaumenze/gjs/releases/"><img src="https://img.shields.io/github/release/chumaumenze/gjs.svg" alt="Latest release"></a>
</p>
</div>


## About

**GJS** provides a comprehensive API to access the WebAssembly host environment when using the js/wasm architecture. Its API is based on JavaScript semantics and provides WebAssembly interop between Go and JS values. Its current scope is to provide a comprehensive and well-tested API.

## Install

```go
GOOS=js GOARCH=wasm go get github.com/chumaumenze/gjs
```

## Usage

```go
package main

import (
	"github.com/chumaumenze/gjs"
)

type Data struct {
	Code    int
	Message string `json:"message"`
	inner   any
}

func main() {
	data := Data{
		Code:    200,
		Message: "Hello World!",
		inner:   "I am ignored",
	}
	_ = gjs.ValueOf(data) // e.g. {Code: 200, "message": "Hello World!"}
}
```

## Motivations

The Go standard library's `syscall/js` package offers limited compatibility between Go and JavaScript values.

- The `js.Value` type lacks support for interfaces or assignment methods for complex Go values.
- The `js.ValueOf` function cannot handle complex types like `struct`s, leading to a panic with the error message `ValueOf: invalid value`.

To overcome these limitations, gjs builds upon and extends the functionality of syscall/js.
