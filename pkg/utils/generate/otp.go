package generate

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func OTP() string {
	sixDigitNum := generateRandomNumber()
	sixDigitStr := formatAsSixDigit(sixDigitNum)
	hash := sha256.Sum256([]byte(sixDigitStr))
	return fmt.Sprintf("%x", hash)
}

func generateRandomNumber() int64 {
	bigInt, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return bigInt.Int64() + 100000
}

func formatAsSixDigit(sixDigitNum int64) string {
	return fmt.Sprintf("%06d", sixDigitNum)
}
