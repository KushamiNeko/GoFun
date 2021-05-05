package hashutils

import (
	crand "crypto/rand"
	"encoding/base64"
	"io"
	mrand "math/rand"
	"net/url"
	"time"
)

func RandBytesB64(length int) string {

	nonce := RandBytes(length)

	return base64.StdEncoding.EncodeToString(nonce)
}

func RandBytes(length int) []byte {
	nonce := make([]byte, length)

	var err error
	for {
		_, err = io.ReadFull(crand.Reader, nonce)
		if err == nil {
			break
		}
	}

	return nonce
}

func RandBytesB64URL(length int) string {
	return url.PathEscape(RandBytesB64(length))
}

func RandString(length int) string {
	const src = "ABCDEFGHIJKLMNOPQLSTUVWXYZabcdefghijklmnopqlstuvwxyz0123456789"

	mrand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = src[mrand.Intn(len(src))]
	}

	return string(b)
}
