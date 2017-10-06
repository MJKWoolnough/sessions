package sessions

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/MJKWoolnough/authenticate"
)

type Store interface {
	Get(*http.Request) []byte
	Set(http.ResponseWriter, []byte)
}

type store struct {
	cookie  http.Cookie
	expires time.Duration
}

func (s *store) GetData(r *http.Request) string {
	for _, cookie := range r.Cookies() {
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

func (s *store) Init(opts ...optFunc) {
	s.cookie.Name = "session"
	s.cookie.Expires = time.Unix(0, 0)
	for _, opt := range opts {
		opt(s)
	}
}

type CookieStore struct {
	store store
	codec authenticate.Codec
}

func NewCookieStore(encKey []byte, opts ...optFunc) (*CookieStore, error) {
	c := new(CookieStore)
	c.store.Init(opts...)
	if cd, err := authenticate.NewCodec(encKey, c.store.expires); err != nil {
		return nil, err
	} else {
		c.codec = *cd
	}
	return c, nil
}

func (c *CookieStore) Get(r *http.Request) []byte {
	data, err := base64.StdEncoding.DecodeString(c.store.GetData(r))
	if err != nil || len(data) < 12 {
		return nil
	}
	dst, _ := c.codec.Decode(data, nil)
	return dst
}

func (c *CookieStore) Set(w http.ResponseWriter, data []byte) {
	if len(data) == 0 {
		c.store.RemoveData(w)
	} else {
		c.store.SetData(w, base64.StdEncoding.EncodeToString(c.codec.Encode(data, nil)))
	}
}

/*
type FSStore struct {
	store
	path string
}

type MemStore struct {
	store
	data sync.Map
}
*/
