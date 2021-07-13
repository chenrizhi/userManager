package main

import (
	"crypto/md5"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	MAXLOGIN = 3
	PASSWORD = "e10adc3949ba59abbe56e057f20f883e"
)

func printMenu() {
	fmt.Println(`
----------------------------------------
1. 查找用户                            |
2. 添加用户                            |
3. 更新用户                            |
4. 删除用户                            |
5. 退出                                |
----------------------------------------`)
}

func login() bool {
	for i := 0; i < MAXLOGIN; i++ {
		fmt.Print("请输入密码：")
		in, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println("输入有误，error:", err)
			continue
		}
		m := fmt.Sprintf("%x", md5.Sum(in))
		if m == PASSWORD {
			fmt.Println("欢迎使用KK用户管理系统！")
			return true
		} else {
			fmt.Println("密码错误。")
		}
	}
	return false
}

func printUser(user user) {
	fmt.Println("id:", user.id)
	fmt.Println("name:", user.name)
	fmt.Println("birthday:", user.birthday.Format("2006-01-02"))
	fmt.Println("tel:", user.tel)
	fmt.Println("address:", user.address)
	fmt.Println("desc:", user.desc)
}

func searchUser(users map[int]*user) {
	fmt.Println("用户查找说明：\n  > 根据用户ID查找：id=1\n  > 根据用户名查找：name=张三\n  > 根据电话号码查找：tel=18888888888")
	in, err := inputString("请输入查找条件：")
	if err != nil {
		fmt.Println("输入错误！")
		return
	}
	inSplit := strings.Split(in, "=")
	if len(inSplit) != 2 {
		fmt.Println("输入格式错误")
		return
	}
	inKey, inValue := strings.TrimSpace(inSplit[0]), strings.TrimSpace(inSplit[1])
	if inKey == "id" {
		inValueInt, err := strconv.Atoi(inValue)
		if err != nil {
			fmt.Println("输入格式错误！error:", err)
			return
		}
		user := users[inValueInt]
		if user != nil {
			printUser(*user)
		}
	} else {
		for _, user := range users {
			inValueInt, err := strconv.Atoi(inValue)
			if err != nil {
				fmt.Println("输入格式错误！error:", err)
				return
			}
			if v := user.id; v == inValueInt {
				printUser(*user)
			}
		}
	}
}

func addUser(users map[int]*user) {
	userId := getId(users)
	name, _ := inputString("姓名：")
	birthdayStr, _ := inputString("生日（2000-01-01）：")
	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		fmt.Println(err)
	}
	tel, _ := inputString("电话：")
	address, _ := inputString("住址：")
	desc, _ := inputString("备注：")
	users[userId] = &user{
		id: userId,
		name: name,
		birthday: birthday,
		tel: tel,
		address: address,
		desc: desc,
	}
}

func updateUser(users map[int]*user) {
	in, _ := inputString("输入要更新的用户ID：")
	inInt, err := strconv.Atoi(in)
	if err != nil {
		fmt.Println("输入错误， error:", err)
		return
	} else if u, ok := users[inInt]; ok {
		name := u.name
		birthday := u.birthday
		tel := u.tel
		address := u.address
		desc := u.desc
		name, _ = inputString("姓名(" + name + ")：")
		birthdayStr, _ := inputString("生日(" + birthday.Format("2006-01-02") + ")：")
		birthday, _ = time.Parse("2006-01-02", birthdayStr)
		tel, _ = inputString("电话(" + tel + ")：")
		address, _ = inputString("住址(" + address + ")：")
		desc, _ = inputString("备注(" + desc + ")：")
		users[inInt] = &user{
			id: inInt,
			name: name,
			birthday: birthday,
			tel: tel,
			address: address,
			desc: desc,
		}
	} else {
		fmt.Println("用户ID不存在")
		return
	}
}

func deleteUser(users map[int]*user) {
	in, _ := inputString("输入要删除的用户ID：")
	if inInt, err := strconv.Atoi(in); err != nil {
		fmt.Println("输入错误， error:", err)
		return
	} else if _, ok := users[inInt]; ok {
		delete(users, inInt)
	} else {
		fmt.Println("用户ID不存在")
		return
	}
}

func inputString(title string) (string, error) {
	var input string
	fmt.Print(title)
	_, err := fmt.Scan(&input)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), err
}

func getId(users map[int]*user) int {
	var max = 1
	for k := range users {
		if k > max {
			max = k
		}
	}
	return max
}

func main() {
	if ! login() {
		return
	}
	users := map[int]*user{}
	opMap := map[string]func(map[int]*user){
		"1": searchUser,
		"2": addUser,
		"3": updateUser,
		"4": deleteUser,
		"5": func(m map[int]*user) {
			os.Exit(0)
		},
	}
	for {
		printMenu()
		op, err := inputString("输入序号：")
		if err != nil {
			fmt.Println("输入错误")
		}
		if opFunc, ok := opMap[op]; ok {
			opFunc(users)
		} else {
			fmt.Println("输入错误")
		}
	}
}
