package util

import (
	"math/rand"
	"strconv"
)

func GenerateRandomEmail() string {
	return "test" + strconv.Itoa(rand.Intn(100000)) + "@localhost"
}
