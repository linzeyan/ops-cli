package password

import (
	"math/rand"
	"strings"
	"time"
)

const (
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbols      = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	numbers      = "0123456789"
	allSet       = lowerLetters + upperLetters + symbols + numbers
)

var p password

type password struct{}

func (password) genString(length uint, charSet string) string {
	rand.Seed(time.Now().Local().UnixNano())
	var s strings.Builder
	for i := uint(0); i < length; i++ {
		s.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return s.String()
}

func GenLower(length uint) string {
	return p.genString(length, lowerLetters)
}

func GenUpper(length uint) string {
	return p.genString(length, upperLetters)
}

func GenSymbol(length uint) string {
	return p.genString(length, symbols)
}

func GenNumber(length uint) string {
	return p.genString(length, numbers)
}

func GenAll(length uint) string {
	return p.genString(length, allSet)
}

func GeneratePassword(length, minLower, minUpper, minSymbol, minNumber uint) string {
	var remain string
	leave := length - minLower - minUpper - minSymbol - minNumber
	lower := GenLower(minLower)
	upper := GenUpper(minUpper)
	symbol := GenSymbol(minSymbol)
	num := GenNumber(minNumber)
	if leave != 0 {
		remain = GenAll(leave)
	}
	result := []byte(lower + upper + symbol + num + remain)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return string(result)
}
