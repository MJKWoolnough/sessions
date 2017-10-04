package sessions

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"sync"
	"time"
)

var timeNow = time.Now

const nonceSize = 12

type codec struct {
	aeadPool sync.Pool
	maxAge   time.Duration
}

func (c *codec) Init(key []byte, maxAge time.Duration) error {
	if l := len(key); l != 16 && l != 24 && l != 32 {
		return errInvalidAES
	}
	a := make([]byte, len(key))
	copy(a, key)
	*c = codec{
		aeadPool: sync.Pool{
			New: func() interface{} {
				block, _ := aes.NewCipher(a)
				aead, _ := cipher.NewGCMWithNonceSize(block, nonceSize)
				return aead
			},
		},
		maxAge: maxAge,
	}
	return nil
}

func (c *codec) Encode(data []byte) string {
	a := c.aeadPool.Get().(cipher.AEAD)

	dst := make([]byte, nonceSize, nonceSize+len(data)+a.Overhead())
	t := timeNow()
	binary.LittleEndian.PutUint64(dst, uint64(t.Nanosecond())) // last four bytes are overriden
	binary.BigEndian.PutUint64(dst[4:], uint64(t.Unix()))

	dst = a.Seal(dst, dst, data, nil)

	c.aeadPool.Put(a)

	return base64.StdEncoding.EncodeToString(dst)
}

func (c *codec) Decode(data string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	} else if len(cipherText) < 12 {
		return nil, errInvalidData
	}

	timestamp := time.Unix(int64(binary.BigEndian.Uint64(cipherText[4:12])), 0)

	if c.maxAge > 0 {
		if t := timeNow().Sub(timestamp); t > c.maxAge || t < 0 {
			return nil, errExpired
		}
	}

	a := c.aeadPool.Get().(cipher.AEAD)

	buf := make([]byte, 0, len(cipherText))

	buf, err = a.Open(buf, cipherText[:12], cipherText[12:], nil)

	c.aeadPool.Put(a)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

var (
	errInvalidAES  = errors.New("invalid AES key, must be 16, 24 or 32 bytes")
	errInvalidData = errors.New("invalid cipher text")
	errExpired     = errors.New("data expired")
)
