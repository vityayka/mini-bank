package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyz"

func RandomMoney() int64 {
	return RandomInt(0, 1_000_000)
}

func RandomName() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.%s", RandomString(6), RandomString(6), RandomString(3))
}

func RandomCurrency() string {
	return GetSupportedCurrencies()[rand.Int31n(3)]
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
