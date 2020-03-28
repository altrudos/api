package charityhonor

import (
	"math/rand"
	"time"

	"github.com/vindexus/randomwords"
)

func init() {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
}

var Generator = randomwords.AdjNounGenerator()

func GenerateUri() string {
	return Generator.Random()
}
