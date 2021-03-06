package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomDouble(min, max float64) float64 {
	return float64(int64((min+rand.Float64()*(max-min+1))*100)) / 100
}

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	sb := strings.Builder{}
	len := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomBalance() float64 {
	return RandomDouble(0, 1000)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currencies := []string{"USD", "RMB", "ERU"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomID() int64 {
	return int64(rand.Intn(10) + 1)
}

func RandomAmount() float64 {
	return RandomDouble(-500, 500)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@emali.com", RandomString(6))
}
