package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/howeyc/gopass"
	"crypto/md5"
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

func printUser(user map[string]string) {
	for k, v := range user {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func searchUser(users map[int]map[string]string) {
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
		printUser(users[inValueInt])
	} else {
		for _, user := range users {
			if v, ok := user[inKey]; ok && v == inValue {
				printUser(user)
			}
		}
	}
}

func addUser(users map[int]map[string]string) {
	userId := getId(users)
	name, _ := inputString("姓名：")
	birthday, _ := inputString("生日（2000-01-01）：")
	tel, _ := inputString("电话：")
	address, _ := inputString("住址：")
	desc, _ := inputString("备注：")
	users[userId] = map[string]string{
		"id": strconv.Itoa(userId),
		"name": name,
		"birthday": birthday,
		"tel": tel,
		"address": address,
		"desc": desc,
	}
}

func updateUser(users map[int]map[string]string) {
	in, _ := inputString("输入要更新的用户ID：")
	if inInt, err := strconv.Atoi(in); err != nil {
		fmt.Println("输入错误， error:", err)
		return
	} else if user, ok := users[inInt]; ok {
		name := user["name"]
		birthday := user["birthday"]
		tel := user["tel"]
		address := user["address"]
		desc := user["desc"]
		name, _ = inputString("姓名(" + name + ")：")
		birthday, _ = inputString("生日(" + birthday + ")：")
		tel, _ = inputString("电话(" + tel + ")：")
		address, _ = inputString("住址(" + address + ")：")
		desc, _ = inputString("备注(" + desc + ")：")
		users[inInt] = map[string]string{
			"id": in,
			"name": name,
			"birthday": birthday,
			"tel": tel,
			"address": address,
			"desc": desc,
		}
	} else {
		fmt.Println("用户ID不存在")
		return
	}
}

func deleteUser(users map[int]map[string]string) {
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

func getId(users map[int]map[string]string) int {
	var max = 1
	for k, _ := range users {
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
	users := map[int]map[string]string{}
	opMap := map[string]func(map[int]map[string]string){
		"1": searchUser,
		"2": addUser,
		"3": updateUser,
		"4": deleteUser,
		"5": func(m map[int]map[string]string) {
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
