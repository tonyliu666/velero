package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	//init function would be called when the first package is first used
	rand.Seed(time.Now().UnixNano()) //Seed function accept int64
	//why rand.Seed? we need different generated data.

}
func RandomInt(min, max int64) int64 {
	//return the value between min and max
	return min + rand.Int63n(max-min+1)
}
func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}
func RandomOwner() string {
	return RandString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
