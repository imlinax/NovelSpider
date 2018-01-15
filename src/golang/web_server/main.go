package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("www/")))
	http.HandleFunc("/api/search", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(404)
			fmt.Println("bad requst")
		}
		err := req.ParseForm()
		if err != nil {
			glog.Info(err)
		}

		novelName := req.FormValue("name")
		fmt.Println("novel name: ", novelName)
		w.Write([]byte(novelName))

	})
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if err != nil {
		glog.Error(err)
	}
}
