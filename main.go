package main

import (
	_ "awesomeProject/userManager/routes"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(http.ListenAndServe(":8080", nil))
}
