package bitesized

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Taken from: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func dasherize(evnt string) string {
	return strings.Join(strings.Split(evnt, " "), "-")
}

// Taken from: http://stackoverflow.com/questions/30272881/unpack-redis-set-bit-string-in-go
func bitStringToBools(str string) []bool {
	bools := make([]bool, 0, len(str)*8)

	for i := 0; i < len(str); i++ {
		for bit := 7; bit >= 0; bit-- {
			isSet := (str[i]>>uint(bit))&1 == 1
			bools = append(bools, isSet)
		}
	}

	return bools
}
