package main

import (
	"awesomeProject/userManager/db"
	_ "awesomeProject/userManager/routes"
	"fmt"
	"net/http"
)

func main() {
	err := db.InitDb("mysql", "golang:golang@tcp(127.0.0.1:3306)/user?charset=utf8mb4&loc=Local&parseTime=true")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(http.ListenAndServe(":8080", nil))
}
