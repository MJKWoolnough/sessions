package sessions

import "time"

type optFunc func(*store)

// Name sets the cookie name.
func Name(name string) optFunc {
	return func(s *store) {
		s.cookie.Name = name
	}
}

// Domain sets the optional domain for the cookie.
func Domain(domain string) optFunc {
	return func(s *store) {
		s.cookie.Domain = domain
	}
}

// Path sets the optional path for the cookie.
func Path(path string) optFunc {
	return func(s *store) {
		s.cookie.Path = path
	}
}

// HTTPOnly sets the httponly flag on the cookie.
func HTTPOnly() optFunc {
	return func(s *store) {
		s.cookie.HttpOnly = true
	}
}

// Secure sets the secure flag on the cookie.
func Secure() optFunc {
	return func(s *store) {
		s.cookie.Secure = true
	}
}

// Expiry sets a maximum time that a cookie and authenticated message are valid
// for.
func Expiry(d time.Duration) optFunc {
	return func(s *store) {
		s.expires = d
	}
}
