package handler

import "net/http"

func Mainpage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main page handler reqsting a method"))
}
