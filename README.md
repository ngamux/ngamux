# ngamux
Simple HTTP router for Go

---

* [Installation](#installation)
* [Examples](#examples)
* [Provided Middlewares](#provided-middlewares)
* [Todo](#todo)

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
    return ngamux.JSON(rw, map[string]string{
      "message": "welcome!",
    })
  })
  
  http.ListenAndServe(":8080", mux)
}
```

# Provided Middlewares
* [Recover](https://github.com/ngamux/ngamux/blob/master/middleware/recover.go)

# Todo
- [x] Multiple handler (middleware for each route)
- [x] Route group
- [x] Route params (in URL parameters)
