package util

import (
	"math/rand"
)

var numAndLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var lowcaseLetters = []rune("abcdefghijklmnopqrstuvwxyz")

var nums = []rune("0123456789")

func randomStr(n int, slice []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = slice[rand.Intn(len(slice))]
	}
	return string(b)
}

func RandNumAndLettersStr(n int) string {
	return randomStr(n, numAndLetters)
}

func RandNum(n int) string {
	return randomStr(n, nums)
}

func RandLowcaseLetters(n int) string {
	return randomStr(n, lowcaseLetters)
}
