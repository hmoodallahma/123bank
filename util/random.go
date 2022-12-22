package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init(){
	// Generate a random int value
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64{
	// generate random int between min and max
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++{
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	stringLen := RandomInt(6,12)
	return RandomString(int(stringLen))
}

func RandomMoney() int64{
	return RandomInt(0,1000)
}


func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}