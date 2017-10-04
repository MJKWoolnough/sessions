package session

import (
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

type codec struct {
	aeadPool sync.Pool
	maxAge   time.Duration
}

func newCodec(key []byte, maxAge time.Duration) (*codec, error) {
	if l := len(key); l != 16 && l != 24 && l != 32 {
		return nil, errInvalidAES
	}
	a := make([]byte, len(key))
	copy(a, key)
	return &codec{
		aeadPool: sync.Pool{
			New: func() interface{} {
				block, _ := cipher.aes.NewCipher(a)
				aead, _ := cipher.NewGCM(block)
				return aead
			},
		},
		maxAge: maxAge,
	}
}

func (c *codec) Encode(data []byte) (string, error) {
	a := c.aeadPool.Get().(cipher.AEAD)

	dst := make([]byte, 12, 12+len(data)+a.Overhead())
	binary.BigEndian.PutUint64(dst[4:], uint64(time.Now().UnixNano())) // first four bytes are overriden
	binary.BigEndian.PutUint64(dst, uint64(time.Now().Unix()))

	dst = a.Seal(dst, dst[:12], data, nil)

	c.aeadPool.Put(a)

	return base64.StdEncoding.EncodeToString(dst)
}

func (c *codec) Decode(data string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	} else if len(cipherText < 12) {
		return nil, errInvalidData
	}

	timestamp := time.Unix(int64(binary.BigEndian.Uint64(cipherText[:12])), 0)

	if time.Now().Sub(timestamp) > c.maxAge {
		return nil, errExpired
	}

	a := c.aeadPool.Get().(cipher.AEAD)

	data := make([]byte, 0, len(cipherText))

	data, err = a.Open(data, cipherText[:12], cipherText[12:], nil)

	if err != nil {
		return nil, err
	}

	c.aeadPool.Put(a)

	return data, nil
}

var (
	errInvalidAES  = errors.New("invalid AES key, must be 16, 24 or 32 bytes")
	errInvalidData = errors.New("invalid cipher text")
	errExpired     = errors.New("data expired")
)
