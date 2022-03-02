# ngamux
Simple HTTP router for Go

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ngamux/ngamux.svg)](https://github.com/ngamux/ngamux)
[![GoDoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/ngamux/ngamux)
[![GoReportCard](https://goreportcard.com/badge/github.com/ngamux/ngamux)](https://goreportcard.com/report/github.com/ngamux/ngamux)
[![Coverage Status](https://coveralls.io/repos/github/ngamux/ngamux/badge.svg?branch=master)](https://coveralls.io/github/ngamux/ngamux?branch=master)

---

* [Installation](#installation)
* [Examples](#examples)
* [Provided Middlewares](#provided-middlewares)
* [License](#license)

---

# Installation
Run this command with correctly configured Go toolchain.
```bash
go get github.com/ngamux/ngamux
```

# Examples
```go
package main

import(
  "net/http"
  "github.com/ngamux/ngamux"
)

func main() {
  mux := ngamux.NewNgamux()
  mux.Get("/", func(rw http.ResponseWriter, r *http.Request) error {
    return ngamux.JSON(rw, ngamux.Map{
      "message": "welcome!",
    })
  })
  
  http.ListenAndServe(":8080", mux)
}
```

# Provided Middlewares
* [CORS](https://github.com/ngamux/middleware/tree/master/cors)
* [Recover](https://github.com/ngamux/middleware/tree/master/recover)
* [Static](https://github.com/ngamux/middleware/tree/master/static)

# License
This project is licensed under the [Mozilla Public License 2.0](https://github.com/ngamux/ngamux/blob/master/LICENSE).

# Contributors
Thanks to all contributors!

[![Contributors](https://contrib.rocks/image?repo=ngamux/ngamux)](https://github.com/ngamux/ngamux/graphs/contributors)
