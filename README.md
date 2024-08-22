# sessions
--
    import "vimagination.zapto.org/sessions"

Package sessions is used to store session information for a web server

## Usage

#### func  Domain

```go
func Domain(domain string) optFunc
```
Domain sets the optional domain for the cookie.

#### func  Expiry

```go
func Expiry(d time.Duration) optFunc
```
Expiry sets a maximum time that a cookie and authenticated message are valid
for.

#### func  HTTPOnly

```go
func HTTPOnly() optFunc
```
HTTPOnly sets the httponly flag on the cookie.

#### func  Name

```go
func Name(name string) optFunc
```
Name sets the cookie name.

#### func  Path

```go
func Path(path string) optFunc
```
Path sets the optional path for the cookie.

#### func  Secure

```go
func Secure() optFunc
```
Secure sets the secure flag on the cookie.

#### type CookieStore

```go
type CookieStore struct {
}
```

CookieStore stores and retrieves authenticated data from a clients cookies.

#### func  NewCookieStore

```go
func NewCookieStore(encKey []byte, opts ...optFunc) (*CookieStore, error)
```
NewCookieStore creates a new CookieStore and initialises it. The options
optFunc's are to set non-default values on the cookie.

#### func (*CookieStore) Get

```go
func (c *CookieStore) Get(r *http.Request) []byte
```
Get retrieves authenticated data from the cookie.

#### func (*CookieStore) Set

```go
func (c *CookieStore) Set(w http.ResponseWriter, data []byte)
```
Set stores authenticated data in a clients cookies.

#### type Store

```go
type Store interface {
	Get(*http.Request) []byte
	Set(http.ResponseWriter, []byte)
}
```

Store is the interface for any stores in this package.
