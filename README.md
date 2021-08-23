# ngamux
Simple HTTP router for Go

---

[Installation](./installation)
[Examples](./examples)

---

# Installation
Run this command with correctly configured Go toolchain.
```
go get github.com/ngamux/ngamux
```

# Examples
```
package main

import(
  "net/http"
  "github.com/ngamux/ngamux"
)

func main() {
  mux := ngamux.NewNgamux()
  mux.Get("/", func(rw http.ResponseWriter, r *http.Request) {
    ngamux.JSON(rw, map[string]string{
      "message": "welcome!",
    })
  })
  
  http.ListenAndServe(":8080", mux)
}
```
