package session

import (
	"net/http"
	"sync"
	"time"
)

type Store interface {
	Get(*http.Request) Session
}

type Session interface {
	Load(k interface{}) (v interface{}, ok bool)
	Store(k, v interface{})
	Delete(k interface{})
}

type session struct {
	mu   sync.RWMutex
	data map[interface{}]interface{}
}

func (s *session) Load(k interface{}) (interface{}, bool) {
	s.mu.RLock()
	v, ok := s.data[k]
	s.mu.RUnlock()
	return v, ok
}

func (s *session) Store(k, v interface{}) {
	s.mu.Lock()
	s.data[k] = v
	s.mu.Unlock()
}

func (s *session) Delete(k interface{}) {
	s.mu.Lock()
	delete(s.data, k)
	s.mu.Unlock()
}

type store struct {
	cookie  http.Cookie
	expires time.Duration
}

func (s *store) getData(r *http.Request) string {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == s.cookie.Name &&
			(s.cookie.Domain == "" || s.cookie.Domain == cookie.Domain) &&
			(s.cookie.Path == "" || s.cookie.Path == cookie.Path) &&
			s.cookie.HttpOnly == cookie.HttpOnly &&
			s.cookie.Secure == cookie.Secure {
			return cookie.Value
		}
	}
	return ""
}

func (s *store) setData(w http.ResponseWriter, data string) {
	cookie := s.cookie
	if s.expires > 0 {
		cookie.Expires = time.Now().Add(s.expires)
	}
	w.Header().Add("Set-Cookie", cookie.String())
}

type optFunc func(s *store)

func newStore(os ...optFunc) *store {
	s := new(store)
	s.cookie.Name = "session"
	for _, o := range os {
		o(s)
	}
	return s
}

type CookieStore struct {
	store
}

type FSStore struct {
	store
	path string
}

type MemStore struct {
	store
	data sync.Map
}
