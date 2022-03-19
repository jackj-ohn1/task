package tools

import (
	"fmt"
	"log"
	"strconv"
	"task/config"
)

func CommonErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

// 为获取cookie单独设置错误，以防陷入获取cookie循环
func GetCookieErr(err error) {
	if err != nil {
		fmt.Println(err, "登录获取cookie失败")
	}
}

func Generate(suffix int, users chan string) {
	var id string
	var grade string
	switch suffix {
	case 21:
		grade = "2021"
	case 20:
		grade = "2020"
	case 19:
		grade = "2019"
	case 18:
		grade = "2018"
	}
	//var num = strconv.I
	for i := 21*10000 + 1; i <= 21*10000+config.Max; i++ {
		id = grade + strconv.Itoa(i)
		users <- id
	}
	fmt.Println("生产完毕!")
	close(users)
}
