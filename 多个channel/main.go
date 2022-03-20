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
	var ch1 = make([]chan string, config.Worker)
	var done = make(chan interface{})
	var now = time.Now()
	// 对每个通道进行初始化
	go func() {
		for i := 0; i < config.Worker; i++ {
			ch1[i] = make(chan string, config.Max/config.Worker)
		}
	}()
	var ch2 = make(chan types.User, 100)

	go tools.Generate(config.Grade, ch1)

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	tools.CommonErr(err)

	client := http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}

	// 初始化client
	time.Sleep(2 * time.Second)

	get.GetCookie(&client, config.Username, config.Password)

	defer close(done) // 不需要缓冲，因为关闭了，那么所有就会停止

	go get.GetInfor(ch1, &client, ch2, done)

	for user := range ch2 {
		count++
		fmt.Println(count, ":", user)
		if count == config.Max {
			close(ch2)
			break
		}
	}

	/* 直接根据woker的完成情况来判断是否结束
	for i := 0; i < len(ch1); i++ {
		<-done
	}
	*/
	fmt.Println("over! time:\n", time.Since(now))
}
