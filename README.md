# sessions

[![CI](https://github.com/MJKWoolnough/sessions/actions/workflows/go-checks.yml/badge.svg)](https://github.com/MJKWoolnough/sessions/actions)
[![Go Reference](https://pkg.go.dev/badge/vimagination.zapto.org/sessions.svg)](https://pkg.go.dev/vimagination.zapto.org/sessions)
[![Go Report Card](https://goreportcard.com/badge/vimagination.zapto.org/sessions)](https://goreportcard.com/report/vimagination.zapto.org/sessions)

--
    import "vimagination.zapto.org/sessions"

Package sessions is used to store session information for a web server.

## Highlights

 - Set and get session data in client cookies.
 - Session data is signed and dated to prevent tampering.

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"vimagination.zapto.org/sessions"
)

func main() {
	store, err := sessions.NewCookieStore([]byte("!THIS IS MY KEY!"), sessions.Expiry(time.Second))
	if err != nil {
		fmt.Println(err)

		return
	}

	w := httptest.NewRecorder()

	store.Set(w, []byte("MY_DATA"))

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	for _, cookie := range w.Result().Cookies() {
		r.AddCookie(cookie)
	}

	fmt.Printf("%q\n", store.Get(r))

	time.Sleep(2 * time.Second)

	fmt.Printf("%q\n", store.Get(r))

	// Output:
	// "MY_DATA"
	// ""
}
```

## Documentation

Full API docs can be found at:

https://pkg.go.dev/vimagination.zapto.org/sessions
