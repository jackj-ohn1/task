package main

import (
	"fmt"
	"strconv"
	"strings"
)

func My_Sprintf(str *string, str1 string, arg [5]int) {
	var index = 100
	str1 = strings.Replace(str1, "%%", "%", -1)
	for i := 0; index >= 0 && i < len(arg); i++ {
		index = strings.Index(str1, "%d")
		if index > 0 {
			front := str1[:index]
			behind := str1[index+2:]
			str1 = front + strconv.Itoa(arg[i]) + behind
			*str = str1
		}
	}
}
func main() {
	var s string
	arg := [5]int{123, 1, 2, 3, 4}
	My_Sprintf(&s, "sfaw%dda%%sa%d", arg)
	fmt.Printf("%s", s)
}
