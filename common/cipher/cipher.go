package cipher

import (
	"math/rand"
	"strings"
	"time"
)

var _chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖabcdefghijklmnopqrstuvwxyzåäö0123456789")

func Random() string {
	rand.Seed(time.Now().UnixNano())
	length := 8 + rand.Intn(8)
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(_chars[rand.Intn(len(_chars))])
	}
	return b.String()
}
