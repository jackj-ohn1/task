package main

import (
	"fmt"
	"log"
	"net/http"
)

type Post struct {
	userName string //用户名
	userPwd  string //密码
	nickname string
	age      string
	sex      string
}

var num = make([]Post, 0, 100)
var infor Post

func register(w http.ResponseWriter, r *http.Request) {
	//获取参数的值，并将其存储在结构体中
	infor.userName = r.FormValue("username")
	infor.userPwd = r.FormValue("password")
	infor.nickname = r.FormValue("nickname")
	infor.age = r.FormValue("age")
	infor.sex = r.FormValue("sex")
	var ok = true
	for _, values := range num {
		if values.userName == infor.userName {
			ok = false
		}
	}
	if ok {
		num = append(num, infor)
		cookies := http.Cookie{
			Name:     infor.userName,
			Value:    infor.userPwd,
			HttpOnly: true,
		}
		w.Header().Set("Set-Cookie", cookies.String())
		w.Write([]byte("注册成功~"))
	}
	if !ok {
		w.Write([]byte("用户已存在~")) //fmt.Fprintln(w,"用户已存在",r.URL.Path),会将“”中的信息和对应的路径后的东西加上去。即：用户已存在 /login
	}
}

func show(p *Post, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
姓名：%s
年龄：%s
密码：%s
性别：%s
昵称：%s
`, p.userName, p.age, p.userPwd, p.sex, p.nickname)
	fmt.Fprintf(w, "\n")
}

func edit(w http.ResponseWriter, r *http.Request) {
	var name = r.FormValue("username")
	var value = r.FormValue("password")
	cookie, err := r.Cookie(name)
	if cookie != nil {
		for i := 0; i < len(num); i++ {
			if name == num[i].userName && num[i].userPwd == value {
				w.Write([]byte("你的信息如下："))
				show(&num[i], w, r)
				num[i].userPwd = r.FormValue("password")
				num[i].nickname = r.FormValue("nickname")
				num[i].age = r.FormValue("age")
				num[i].sex = r.FormValue("sex")
				w.Write([]byte("修改后的信息如下：\n"))
				show(&num[i], w, r)
			}
		}
	} else {
		fmt.Fprintln(w, err)
	}
}

func show_all(w http.ResponseWriter, r *http.Request) {
	var tmp = r.FormValue("username")
	cookie, err := r.Cookie(tmp)
	if cookie != nil {
		if cookie.Name == tmp && cookie.Value == r.FormValue("password") {
			for index, value := range num {
				fmt.Fprintf(w, `已注册的第%d的人的信息:
姓名：%s
年龄：%s
密码：%s
性别：%s
昵称：%s
`, index+1, value.userName, value.age, value.userPwd, value.sex, value.nickname)
				fmt.Fprintf(w, "\n")
			}
		}
	} else {
		fmt.Fprintln(w, err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	for _, value := range num {
		if value.userName == r.FormValue("username") && value.userPwd == r.FormValue("password") {
			cookies := http.Cookie{
				Name:     infor.userName,
				Value:    infor.userPwd,
				HttpOnly: true,
			}
			w.Header().Set("Set-Cookie", cookies.String())
			w.Write([]byte("登陆成功！"))
		}
		if value.userName == r.FormValue("username") && value.userPwd != r.FormValue("password") {
			w.Write([]byte("密码错误！"))
		}
	}
}

func main() {
	//对  /login 进行reigister函数操作
	http.HandleFunc("/register", register)
	http.HandleFunc("/register/edit", edit)
	http.HandleFunc("/register/show", show_all)
	http.HandleFunc("/login", login)
	http.HandleFunc("/login/edit", edit)
	http.HandleFunc("/login/show", show_all)

	//创建监听的路由端口
	err := http.ListenAndServe(":6789", nil)
	if err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
