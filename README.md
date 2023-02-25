# Autoreleasepool memory management for golang

![License](https://img.shields.io/github/license/demdxx/autoreleasepool)
[![GoDoc](https://godoc.org/github.com/demdxx/autoreleasepool?status.svg)](https://godoc.org/github.com/demdxx/autoreleasepool)
[![Testing Status](https://github.com/demdxx/autoreleasepool/workflows/Tests/badge.svg)](https://github.com/demdxx/autoreleasepool/actions?workflow=Tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/demdxx/autoreleasepool)](https://goreportcard.com/report/github.com/demdxx/autoreleasepool)
[![Coverage Status](https://coveralls.io/repos/github/demdxx/autoreleasepool/badge.svg?branch=main)](https://coveralls.io/github/demdxx/autoreleasepool?branch=main)

> Minimal version of golang 1.20

Most of the golang applications are web services, which are processing a lot of requests.
In this case, the memory management is very important. The main problem is that the golang
does not have a garbage collector, which can free the memory after the request is processed.
The main idea of this package is to create a pool of autoreleasepools, which can be used
to free the memory after the request is processed.

This package has been inspired by [@autoreleasepool](https://developer.apple.com/library/archive/documentation/Cocoa/Conceptual/MemoryMgmt/Articles/mmAutoreleasePools.html)
of Objective-C and Swift and can be used in the same way.
All memory allocated in the autoreleasepool will be freed after the request is processed,
wich can significantly speed up the application and reduce the garbage collector load.

The `"arena"` package is experimental and must be activated by `GOEXPERIMENT=arenas`.

## Example

```go
package main

import (
  "fmt"
  "context"

  "github.com/demdxx/autoreleasepool/httpautoreleasepool"
  "github.com/demdxx/autoreleasepool"
)

type RenderContext struct {
  name string
  age string
}

func indexHttpHandler(w http.ResponseWriter, r *http.Request) {
  ct := autoreleasepool.New[RenderContext](r.Context())
  ct.name = "Dem"
  ct.age = "30"
  fmt.Fprintf(w, "Hello, %s! You are %s years old.", ct.name, ct.age)
}

func main() {
  http.HandleFunc("/", httpautoreleasepool.Wrap(indexHttpHandler))
  http.ListenAndServe(":8080", nil)
}
```