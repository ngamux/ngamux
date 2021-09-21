# ngamux
Simple HTTP router for Go

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
