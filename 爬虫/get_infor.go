package main

import (
	"fmt"
	"io/ioutil"
	"mygo/myself"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Trans(a string) string {
	b := strconv.QuoteToASCII(a)
	c := b[1 : len(b)-1]
	d := strings.Split(c, "\\u")
	var e string
	for _, v := range d {
		if len(v) < 1 {
			continue
		}
		tmp, _ := strconv.ParseInt(v, 16, 32)

		e += fmt.Sprintf("%c", tmp)
	}
	return e
}

func main() {
	file, err := os.Create("C:\\source\\content.txt")
	if err != nil {
		fmt.Println("打开失败！", err)
	} else {
		fmt.Println("打开成功！")
	}
	key := myself.Get()
	for i := 1; i <= 5; i++ {
		url := "http://work.muxi-tech.xyz/api/v1.0/status/list/" + strconv.Itoa(i) + "/"

		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		Error(err)
		req.Header.Add("token", key)

		res, err := client.Do(req)
		Error(err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		Error(err)
		//使用正则来匹配对应的内容

		//获取name
		re1 := regexp.MustCompile(`"username": "(.*?)"`)
		result1 := re1.FindAllString(string(body), -1)

		//获取发表时间
		re2 := regexp.MustCompile(`"time": "(.*?)"`)
		result2 := re2.FindAllString(string(body), -1)
		//获取内容
		re3 := regexp.MustCompile(`"content": "(.*?)"`)
		result3 := re3.FindAllString(string(body), -1)

		re := regexp.MustCompile(`\\[a-z0-9]{5}`)

		for i := 0; i < len(result1); i++ { //单纯去统一一个数
			//找到每一个人对应的数据
			content1 := re.FindAllString(result1[i], -1)
			content2 := result2[i]
			content3 := re.FindAllString(result3[i], -1)
			var name string
			var str string

			//将每一个人对应的姓名，发表时间，内容写入文件
			for j := len(content1) - 1; j >= 0; j-- {
				name = Trans(content1[j]) + name
			}
			_, _ = file.Write([]byte(`
		`)) //提高可观性
			a, err := file.Write([]byte(strconv.Itoa(i+1) + ".姓名：" + name)) //写入时，是从 name的后面往前读取，所以把循环的方向改了
			if err != nil {
				fmt.Println("写入失败！", a, err)
			} else {
				fmt.Println("写入成功！", a)
			}

			_, _ = file.Write([]byte(`
		`)) // time无需做处理

			a, err = file.Write([]byte(content2))
			if err != nil {
				fmt.Println("写入失败！", a, err)
			} else {
				fmt.Println("写入成功！", a)
			}

			for j := len(content3) - 1; j >= 0; j-- {
				str = Trans(content3[j]) + str
			}
			_, _ = file.Write([]byte(`
		`))
			a, err = file.Write([]byte("内容：" + str))
			if err != nil {
				fmt.Println("写入失败！", a, err)
			} else {
				fmt.Println("写入成功！", a)
			}
			_, _ = file.Write([]byte(`
		
		`))
		}
	}
}
