package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	const charset = CharsetNumeric
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range result {
		num, err := rand.Int(rand.Reader, charsetLen)

		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result),nil
}

func GenerateInviteCode(length int) (string, error) {
	const safeCharset = CharsetSafe
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(safeCharset))))
		if err != nil {
			return "", err
		}
		result[i] = safeCharset[num.Int64()]
	}

	return string(result), nil
}
