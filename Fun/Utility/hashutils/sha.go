package hashutils

import (
	"encoding/base64"

	"golang.org/x/crypto/sha3"
)

func ShaB64FromString(srcs ...string) string {
	bytesSlice := make([][]byte, len(srcs))

	for i, src := range srcs {
		bytesSlice[i] = []byte(src)
	}

	cipher := ShaB64FromBytes(bytesSlice...)

	return cipher
}

func ShaB64FromBytes(srcs ...[]byte) string {

	if len(srcs) == 0 {
		return ""
	}

	h := sha3.New512()

	for _, src := range srcs {
		h.Write(src)
	}

	hash := h.Sum(nil)

	cipher := base64.StdEncoding.EncodeToString(hash)

	return cipher
}
