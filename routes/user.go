package routes

import (
	"awesomeProject/userManager/controllers"
	"net/http"
)

func init()  {
	http.HandleFunc("/", controllers.IndexAction)
	http.HandleFunc("/users/add/", controllers.AddAction)
	http.HandleFunc("/users/modify/", controllers.ModifyAction)
	http.HandleFunc("/users/delete/", controllers.DeleteAction)
}