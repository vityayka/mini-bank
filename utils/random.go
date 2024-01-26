package utils

import (
	"math/rand"
	"strings"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyz"

func RandomMoney() int64 {
	return randomInt(0, 1_000_000)
}

func RandomName() string {
	return randomString(6)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "UAH"}
	return currencies[rand.Int31n(3)]
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
