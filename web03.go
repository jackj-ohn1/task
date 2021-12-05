package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func Error(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var db, err = sql.Open("mysql", "yyj:0118@(127.0.0.1)/FIRST") //用户名:密码@ip地址/要链接的数据库的名字

func main() {
	if err != nil {
		fmt.Println("连接失败", err)
	} else {
		fmt.Println("连接成功")
	}
	defer db.Close()
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/login/edit", Edit)
	http.HandleFunc("/login/view", View)
	http.HandleFunc("/login/delete", Delete)
	err = http.ListenAndServe("localhost:1234", nil)
	Error(err)
}

//注册
func Register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")
	tmp1 := r.FormValue("age")
	sex := r.FormValue("sex")
	id := r.FormValue("id")
	age, err := strconv.Atoi(tmp1)
	Error(err)
	var (
		id_tmp string
	)
	sql := "SELECT id FROM user"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()
	var ok = true
	for rows.Next() {
		err = rows.Scan(&id_tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
		if id_tmp == id {
			ok = false
		}
	}
	if ok {
		Insert(name, password, age, sex, id)
		cookie := http.Cookie{
			Name:  "cookie",
			Value: id,
		}
		Set(cookie, w, r)
		fmt.Fprintln(w, "注册成功！")
	} else {
		fmt.Fprintln(w, "用户已存在")
	}
}

//登录
func Login(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	id := r.FormValue("id")
	var (
		id_tmp       string
		password_tmp string
	)
	sql := "SELECT id,password FROM user"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()
	var ok = true
	for rows.Next() {
		err = rows.Scan(&id_tmp, &password_tmp)
		if err != nil {
			fmt.Println(err)
			return
		}
		if id_tmp == id && password == password_tmp {
			cookie := http.Cookie{
				Name:  "cookie", //用id 来辨识一个人的身份
				Value: id,
			}
			_, err := r.Cookie("cookie")
			if err != nil {
				Set(cookie, w, r)
			} else {
				fmt.Fprintln(w, "已存有用户信息，请退出重试")
				break
			}
			ok = false
			fmt.Fprintln(w, "请输入辨识身份的name，id和password！")
			fmt.Fprint(w, "登陆成功！")
			break
		}
	}
	if ok {
		fmt.Fprintln(w, "您输入的信息有误，请重新检查！")
	}
}

//查询数据
func View(w http.ResponseWriter, r *http.Request) {
	var (
		name     string
		password string
		age      int
		sex      string
		id       string
		count    int
	)
	sql := "SELECT * FROM user"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()
	for rows.Next() {
		count++
		err = rows.Scan(&name, &password, &age, &sex, &id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, "姓名：%s 密码：%s 年龄：%d 性别：%s ID：%s\n", name, password, age, sex, id)
	}
	fmt.Fprintf(w, "共有%d人已注册。", count)
}

//对学生信息进行编辑
func Edit(w http.ResponseWriter, r *http.Request) {
	id := Get(w, r)
	_, ok := Look(id)
	if ok {
		Update(w, r, id)
	}
}

//查询部分数据
func Look(id string) (map[string]string, bool) {
	var (
		id_tmp   string
		name     string
		age      int
		sex      string
		password string
	)
	ok := false
	tmp := make(map[string]string, 4)
	sql := "SELECT name,password,age,sex,id FROM user"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("查找失败", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&name, &password, &age, &sex, &id_tmp)
		if err != nil {
			fmt.Println(err)
		}
		if id == id_tmp {
			tmp["name"] = name
			tmp["age"] = strconv.Itoa(age)
			tmp["sex"] = sex
			tmp["password"] = password
			ok = true
		}
	}
	return tmp, ok
}

//添加数据
func Insert(name string, password string, age int, sex string, id string) {
	fmt.Println("开始插入数据")
	sql := "INSERT INTO user(name,password,age,sex,id) VALUES(?,?,?,?,?)"
	_, err := db.Exec(sql, name, password, age, sex, id)
	if err != nil {
		fmt.Println("插入失败", err)
	} else {
		fmt.Println("插入成功")
	}
}

//删除学生
func Delete(w http.ResponseWriter, r *http.Request) {
	sql := "DELETE FROM user WHERE id=?" //指定被删除者的信息
	id := r.FormValue("id")
	_, err := db.Exec(sql, id)
	if err != nil {
		fmt.Println("删除失败", err)
	} else {
		fmt.Println("删除成功")
	}
	w.Write([]byte("请输入将被删除的学生的id"))
}

//更新数据

func Update(w http.ResponseWriter, r *http.Request, id string) {
	//在未输入某个字段时，使其值不发生改变
	var tmp int
	fmt.Fprintln(w, "请输入需要改变后的信息！")
	name := r.FormValue("name")
	password := r.FormValue("password")
	sex := r.FormValue("sex")
	age := r.FormValue("age")
	//age, _ := strconv.Atoi(tmp)
	sql := "UPDATE user SET name=?,password=?,sex=?,age=? WHERE id=?"
	infor, _ := Look(id)
	fmt.Println(infor)
	fmt.Println(password, name, sex, age)
	if password == "" {
		name = infor["password"]
	}
	if sex == "" {
		sex = infor["sex"]
	}
	if age == "" {
		age = infor["age"]
		tmp, _ = strconv.Atoi(age)
	} else {
		tmp, _ = strconv.Atoi(age)
	}
	if name == "" {
		name = infor["name"]
	}

	_, err := db.Exec(sql, name, password, sex, tmp, id)
	if err != nil {
		fmt.Fprintln(w, "修改失败", err)
	} else {
		fmt.Fprintln(w, "修改成功")
	}
}

//发送cookie
func Set(cookie http.Cookie, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Set-Cookie", cookie.String())
}

//获取cookie
func Get(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("cookie")
	if err != nil {
		fmt.Println("获取cookie失败！", err)
	}
	return cookie.Value
}
