package main

import (
	"net/http"

	"github.com/golang/glog"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("www/")))
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if err != nil {
		glog.Error(err)
	}
}
