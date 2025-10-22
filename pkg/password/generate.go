package password

import (
	"math/rand"
	"strings"
)

func Generate(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!@.*&%$#"

	bldr := strings.Builder{}

	for i := 0; i < length; i++ {
		ind := rand.Intn(len(charset))
		bldr.WriteByte(charset[ind])
	}

	return bldr.String()
}
