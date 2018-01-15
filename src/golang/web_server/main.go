package main

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		data, err := ioutil.ReadFile("www/index.html")
		if err != nil {
			glog.Error(err)
		}
		w.Write(data)
	})
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if err != nil {
		glog.Error(err)
	}
}
