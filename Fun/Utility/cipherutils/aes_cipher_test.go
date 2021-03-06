package cipherutils

import (
	"strings"
	"testing"
)

func TestAesCipherConsistency(t *testing.T) {

	t.Parallel()

	secret := "lalalalalala~~~~~"

	cipherText := AesCipherInstanceEncode(secret)
	plainText, err := AesCipherInstanceDecode(cipherText)
	if err != nil {
		t.Errorf("\nError: %v\n", err)
	}

	if strings.Compare(secret, plainText) != 0 {
		t.Errorf("\nExpect: %v\nGet: %v\n", secret, plainText)
	}

}
