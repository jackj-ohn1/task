package get

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"task/config"
	"task/pkg/tools"
	"task/pkg/types"
	"time"

	"golang.org/x/net/publicsuffix"
)

// 获取信息
func GetInfor(ch1 []chan string, client *http.Client, ch2 chan<- types.User, done chan interface{}) {
	var user types.User

	for i := 0; i < len(ch1); i++ {
		fmt.Printf("信息正在获取中...(work %d)\n", i+1)
		// 否则会发生闭包
		go func(i int) {
			for {
				select {
				case id, ok := <-ch1[i]:
					if ok {
						url := "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term=" + id + "&_=1647259625553"
						req, err := http.NewRequest("GET", url, nil)

						tools.CommonErr(err)

						resp, err := client.Do(req)

						clienterr := CookieErr(err)

						// 防止爬取时，cookie失效
						if clienterr != nil {
							resp, _ = clienterr.Do(req)
						}

						body, err := io.ReadAll(resp.Body)

						tools.CommonErr(err)

						name_reg := regexp.MustCompile(`"name": "(.*?)"`)
						re := name_reg.FindAllStringSubmatch(string(body), -1)

						user.Id = id

						if len(re) == 0 {
							user.Name = "查无此人"
						} else {
							user.Name = re[0][1]
						}

						ch2 <- user
						time.Sleep(200 * time.Millisecond) // 适当停止一下，防止访问次数过多导致要重新获取
					}
				case <-done:
					fmt.Printf("work %d死亡\n", i+1)
					return
				case <-time.After(1 * time.Second):
					fmt.Printf("work %d非正常死亡!\n", i+1)
					return
				}
			}
		}(i)
	}
}

// 获取登录需要上传的表单值
func GetbaseForm(client *http.Client) types.Form {

	var form types.Form

	req, err := http.NewRequest("GET", "http://kjyy.ccnu.edu.cn/loginall.aspx", nil)

	tools.CommonErr(err)

	resp, err := client.Do(req)

	tools.CommonErr(err)

	body, err := io.ReadAll(resp.Body)

	tools.CommonErr(err)

	// 获取 hidden 对应的表单值
	lt_reg := regexp.MustCompile(`name="lt" value="(.*?)"`)
	_eventId_reg := regexp.MustCompile(`name="_eventId" value="(.*?)"`)
	execution_reg := regexp.MustCompile(`name="execution" value="(.*?)"`)

	form.Lt = lt_reg.FindAllStringSubmatch(string(body), -1)[0][1]
	form.EventId = _eventId_reg.FindAllStringSubmatch(string(body), -1)[0][1]
	form.Execution = execution_reg.FindAllStringSubmatch(string(body), -1)[0][1]

	/*for _, v := range resp.Cookies() {
		jsessionid = v.Value
	}*/
	fmt.Println("表单获取成功!")
	return form
}

// 登陆成功就把cookie存放到cookiejar里面
func GetCookie(client *http.Client, username, password string) {

	var user = url.Values{}

	form := GetbaseForm(client)

	user.Set("execution", form.Execution)
	user.Set("_eventId", form.EventId)
	user.Set("lt", form.Lt)
	user.Set("password", password)
	user.Set("username", username)

	req, err := http.NewRequest("POST", "https://account.ccnu.edu.cn/cas/login?service=http://kjyy.ccnu.edu.cn/loginall.aspx?page=", strings.NewReader(user.Encode()))

	tools.CommonErr(err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Mobile Safari/537.36 Edg/99.0.1150.39")

	_, err = client.Do(req)

	tools.CommonErr(err)

	fmt.Println("cookie获取成功!")

}

// 如果cookie发生错误，很可能是过期了所以我们就重新获取
func CookieErr(err error) *http.Client {
	jar, errnew := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	tools.CommonErr(errnew)

	client := http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
	}
	if err != nil {
		GetCookie(&client, config.Username, config.Password)
		//log.Println("err:", err)
		return &client
	}
	return nil
}
