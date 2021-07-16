package controllers

import (
	"awesomeProject/userManager/user"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func IndexAction(writer http.ResponseWriter, request *http.Request) {
	tpl := template.Must(template.New("index.tpl.html").ParseFiles("templates/index.tpl.html"))
	search := strings.TrimSpace(request.FormValue("search"))
	var err error

	if search == "" {
		err = tpl.Execute(writer, user.GetUsers())
	} else {
		err = tpl.Execute(writer, user.GetUsers(search))
	}

	if err != nil {
		fmt.Println("error:", err)
	}
}

func AddAction(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		fmt.Println("error:", err)
		return
	}
	u := user.User{
		Name:     request.FormValue("name"),
		Address:  request.FormValue("address"),
		Birthday: request.FormValue("birthday"),
		Tel:      request.FormValue("tel"),
		Remarks:  request.FormValue("remarks"),
		Password: request.FormValue("password"),
	}
	tpl := template.Must(template.New("add.tpl.html").ParseFiles("templates/add.tpl.html"))
	if http.MethodGet == request.Method {
		err := tpl.Execute(writer, u)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else if http.MethodPost == request.Method {
		user.InsertUser(u)
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func ModifyAction(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		fmt.Println("error:", err)
		return
	}
	userId, err := strconv.Atoi(request.FormValue("id"))
	if err != nil {
		fmt.Println("modify user conv userId failed:", err)
		return
	}
	u := user.GetUser(userId)
	tpl := template.Must(template.New("modify.tpl.html").ParseFiles("templates/modify.tpl.html"))
	if http.MethodGet == request.Method {
		err := tpl.Execute(writer, u)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else if http.MethodPost == request.Method {
		u := user.GetUser(userId)
		u.Name = request.FormValue("name")
		u.Address = request.FormValue("address")
		u.Birthday = request.FormValue("birthday")
		u.Tel = request.FormValue("tel")
		u.Remarks = request.FormValue("remarks")
		u.Password = request.FormValue("password")
		user.UpdateUser(*u)
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func DeleteAction(writer http.ResponseWriter, request *http.Request) {
	userId, err := strconv.Atoi(request.FormValue("id"))
	if err != nil {
		fmt.Println("userId conv failed:", err)
		http.Redirect(writer, request, "/", http.StatusFound)
	}
	err = user.DeleteUser(userId)
	if err != nil {
		fmt.Println("delete user failed:", err)
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}