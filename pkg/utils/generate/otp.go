package generate

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const (
	otpMinValue = 100000
	otpRange    = 900000
)

func OTP() string {
	sixDigitNum := generateRandomNumber()
	sixDigitStr := formatAsSixDigit(sixDigitNum)
	hash := sha256.Sum256([]byte(sixDigitStr))
	return fmt.Sprintf("%x", hash)
}

func generateRandomNumber() int64 {
	bigInt, _ := rand.Int(rand.Reader, big.NewInt(otpRange))
	return bigInt.Int64() + otpMinValue
}

func formatAsSixDigit(sixDigitNum int64) string {
	return fmt.Sprintf("%06d", sixDigitNum)
}
