package sessions

import "time"

type optFunc func(*store)

func Name(name string) optFunc {
	return func(s *store) {
		s.cookie.Name = name
	}
}

func Domain(domain string) optFunc {
	return func(s *store) {
		s.cookie.Domain = domain
	}
}

func Path(path string) optFunc {
	return func(s *store) {
		s.cookie.Path = path
	}
}

func HTTPOnly() optFunc {
	return func(s *store) {
		s.cookie.HttpOnly = true
	}
}

func Secure() optFunc {
	return func(s *store) {
		s.cookie.Secure = true
	}
}

func Expiry(d time.Duration) optFunc {
	return func(s *store) {
		s.expires = d
	}
}
