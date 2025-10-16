package sessions_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"vimagination.zapto.org/sessions"
)

func Example() {
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
