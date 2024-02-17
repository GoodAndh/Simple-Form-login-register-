package main

import (

	jalur "webpart2/path"
	"net/http"
)

func main() {
	http.HandleFunc("/", jalur.Index)
	http.HandleFunc("/login", jalur.Login)
	http.HandleFunc("/logout", jalur.Logout)
	http.HandleFunc("/register", jalur.Register)

	http.ListenAndServe(":8080",nil)
}
