package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"runtime"
	"task/config"
	"task/pkg/get"
	"task/pkg/tools"
	"task/pkg/types"
	"time"

	"golang.org/x/net/publicsuffix"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var count = 0
	var ch1 = make(chan string, 100)
	var ch2 = make(chan types.User, 100)
	var done = make(chan interface{})
	var now = time.Now()

	go tools.Generate(config.Grade, ch1)

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	tools.CommonErr(err)

	client := http.Client{
		Timeout: 5 * time.Second,
		Jar:     jar,
	}

	// 初始化client

	get.GetCookie(&client, config.Username, config.Password)

	defer close(done) // 不需要缓冲，因为关闭了，那么所有就会停止

	for i := 1; i <= config.Worker; i++ { // 额，只用一个还更快，因为多个使用的话更可能会出现访问错误，进而导致重新获取cookie
		go get.GetInfor(i, ch1, &client, ch2, done)
	}

	for user := range ch2 {
		count++
		fmt.Println(count, ":", user)
		if count == config.Max {
			close(ch2)
			break
		}
	}
	fmt.Println("over! time:\n", time.Since(now))
}
