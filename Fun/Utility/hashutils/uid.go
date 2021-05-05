package hashutils

import "os"

func GenerateUID(user, password string) string {
	UID := ShaB64FromString(user, password, os.Getenv("UID_SECRET"))

	return UID
}
