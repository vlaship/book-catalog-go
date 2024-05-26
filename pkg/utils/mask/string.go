package mask

const (
	maskedPart = "*****"
	minLen     = 5
	cut        = 2
)

func String(s string) string {
	if len(s) <= minLen {
		return maskedPart
	}
	return s[:cut] + maskedPart + s[len(s)-cut:]
}
