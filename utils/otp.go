package utils

import "math/rand"

func GenerateOtp(otp int) string {
	var setOtp = []rune("0123456789")
	number := make([]rune, otp)
	for i := range number {
		number[i] = setOtp[rand.Intn(len(setOtp))]
	}
	return string(number)
}