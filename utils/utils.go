package utils

import (
	"math/rand"
	"time"
)

func Token(l int) string {

	str := "0123456789abcdefhijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	result := []byte{}
	for i := 0; i < l; i++ {
		result = append(result, str[rand.Intn(len(str))])
	}
	return string(result)
}

func Finger(len int, prefix string) string {
	var str_li = []string{}
	for i := 0; i < len; i++ {

		str := Token(len)
		time.Sleep(time.Nanosecond * 10)
		str_li = append(str_li, str)
	}
	finger := prefix + "-" + str_li[0] + "-" + str_li[1] + "-" + str_li[2]
	return finger
}
