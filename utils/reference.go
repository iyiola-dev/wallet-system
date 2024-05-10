package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func GenerateReference() (string, error) {
	//return random strings of 6 characters
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, 11)
	for i := 0; i < 11; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", errors.New("error generating reference")
		}
		ret[i] = letters[num.Int64()]
	}
	reference := strings.ToUpper(fmt.Sprintf("%s", string(ret)))

	return reference, nil
}
