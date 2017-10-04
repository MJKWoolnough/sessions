package session

import (
	"net/http"
	"sync"
	"time"
)

type Store interface {
	Get(*http.Request) []byte
}

type store struct {
	cookie  http.Cookie
	expires time.Duration
}

func (s *store) GetData(r *http.Request) string {
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

func (s *store) SetData(w http.ResponseWriter, data string) {
	cookie := s.cookie
	if s.expires > 0 {
		cookie.Expires = time.Now().Add(s.expires)
	}
	cookie.Value = data
	w.Header().Add("Set-Cookie", cookie.String())
}

func (s *store) RemoveData(w http.ResponseWriter) {
	w.Header().Add("Set-Cookie", s.cookie.String())
}

type optFunc func(s *store)

func (s *store) Init(opts ...optFunc) {
	s.cookie.Name = "session"
	for _, opt := range opts {
		opt(s)
	}
}

type CookieStore struct {
	store store
	codec codec
}

func NewCookieStore(encKey []byte, opts ...optFunc) (*CookieStore, error) {
	c := new(CookieStore)
	c.store.Init(opts...)
	if err := c.codec.Init(encKey, c.store.expires); err != nil {
		return nil, err
	}
	return c, nil
}

type FSStore struct {
	store
	path string
}

type MemStore struct {
	store
	data sync.Map
}
