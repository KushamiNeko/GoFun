package cipherutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"

	"github.com/KushamiNeko/GoFun/Utility/hashutils"
)

const (
	storeCapacity = 9
)

type aesCipherInstance struct {
	Key   []byte
	Nonce []byte

	ChipherBlock cipher.Block

	ChipherGCM cipher.AEAD
}

var (
	aesCipherInstanceStore []*aesCipherInstance
)

func init() {
	err := AesCipherInstanceInit()
	if err != nil {
		panic(err)
	}
}

func AesCipherInstanceInit(args ...interface{}) error {

	aesCipherInstanceStore = make([]*aesCipherInstance, storeCapacity)
	var err error

	for i := 0; i < storeCapacity; i++ {

		instance := new(aesCipherInstance)
		instance.Key = hashutils.RandBytes(32)
		instance.Nonce = hashutils.RandBytes(12)

		instance.ChipherBlock, err = aes.NewCipher(instance.Key)
		if err != nil {
			return err
		}

		instance.ChipherGCM, err = cipher.NewGCM(instance.ChipherBlock)
		if err != nil {
			return err
		}

		aesCipherInstanceStore[i] = instance
	}

	return nil

}

func AesCipherInstanceEncode(plainTextString string) string {

	instance := aesCipherInstanceStore[rand.Intn(storeCapacity)]

	plainText := []byte(plainTextString)

	cipherText := instance.ChipherGCM.Seal(nil, instance.Nonce, plainText, nil)

	cipherTextString := base64.StdEncoding.EncodeToString(cipherText)

	return cipherTextString
}

func AesCipherInstanceDecode(cipherTextString string) (string, error) {

	var plainText []byte
	var err error

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextString)
	if err != nil {
		return "", err
	}

	for _, v := range aesCipherInstanceStore {
		plainText, err = v.ChipherGCM.Open(nil, v.Nonce, cipherText, nil)
		if err == nil {
			return string(plainText), nil
		}
	}

	return "", err
}
