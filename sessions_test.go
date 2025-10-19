package sessions

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	store, err := NewCookieStore([]byte("!THIS IS MY KEY!"), Expiry(time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	w := httptest.NewRecorder()

	const myData = "MY_DATA"

	store.Set(w, []byte(myData))

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	for _, cookie := range w.Result().Cookies() {
		r.AddCookie(cookie)
	}

	if data := store.Get(r); string(data) != myData {
		t.Errorf("expected to read %q, got %q", myData, data)
	}

	time.Sleep(2 * time.Second)

	if data := store.Get(r); data != nil {
		t.Errorf("expected to read nil, got %q", data)
	}
}

func TestSign(t *testing.T) {
	t.Parallel()

	store, err := NewCookieStore([]byte("!THIS IS MY KEY!"), Expiry(time.Second), Sign())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	w := httptest.NewRecorder()

	const myData = "MY_DATA"

	store.Set(w, []byte(myData))

	r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	for _, cookie := range w.Result().Cookies() {
		r.AddCookie(cookie)
	}

	if data := store.Get(r); string(data) != myData {
		t.Errorf("expected to read %q, got %q", myData, data)
	}

	time.Sleep(2 * time.Second)

	if data := store.Get(r); data != nil {
		t.Errorf("expected to read nil, got %q", data)
	}
}
