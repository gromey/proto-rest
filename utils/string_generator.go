package utils

import (
	"math/rand"
	"time"
)

var (
	Nums       characterSet = []rune("0123456789")
	UpLetters  characterSet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	LowLetters characterSet = []rune("abcdefghijklmnopqrstuvwxyz")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// characterSet a set of characters to generate a random string.
type characterSet []rune

func (cs characterSet) Append(sets ...[]rune) []rune {
	length := len(cs)
	for _, set := range sets {
		length += len(set)
	}

	newSet := make([]rune, length)

	i := copy(newSet, cs)

	for _, set := range sets {
		i += copy(newSet[i:], set)
	}

	return newSet
}

type RandomString struct {
	prefix       string
	characterSet []rune
	length       int
}

// Configure sets configuration for generates random string.
func (s *RandomString) Configure(prefix string, charSet []rune, randLength int) {
	s.prefix = prefix
	s.characterSet = charSet
	s.length = randLength
}

// Generate generates random string according to configuration.
func (s *RandomString) Generate() string {
	b := make([]rune, s.length)

	for i := range b {
		b[i] = s.characterSet[rand.Intn(len(s.characterSet))]
	}

	return s.prefix + string(b)
}
