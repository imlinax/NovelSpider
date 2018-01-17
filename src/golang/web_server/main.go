package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"
)

const (
	TOOLDIR = "tools"
)

func main() {
	flag.Parse()
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
		w.Write([]byte("starting get " + novelName))
		go crawlNovel(novelName)

	})
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if err != nil {
		glog.Error(err)
	}
}

func crawlNovel(name string) {
	crawlerPath := filepath.Join(TOOLDIR, "crawl_novel")
	cmd := exec.Command(crawlerPath, "-name="+name)
	err := cmd.Run()
	if err != nil {
		glog.Error(err)
	}
	glog.Info("下载完成: ", name)
}
