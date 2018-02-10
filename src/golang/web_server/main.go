package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"

	"github.com/golang/glog"
)

const (
	SEARCH_SITE = `http://www.biquge.com.tw`
	SEARCH_URL  = SEARCH_SITE + `/modules/article/soshu.php?searchkey=%s`
	TOOLDIR     = "tools"
)

type NovelInfo struct {
	Name          string
	Author        string
	LatestChapter string
	EntryLink     string
	DownloadLink  string
	SubScribeLink string
}

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.Redirect(w, req, "/www/index.html", http.StatusFound)
	} else {
		http.Redirect(w, req, "/www/404.html", http.StatusFound)
	}

}
func searchHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		http.Redirect(w, req, "/404.html", 404)
		fmt.Println("bad requst")
		return
	}
	err := req.ParseForm()
	if err != nil {
		glog.Info(err)
	}

	novelName := req.FormValue("name")
	fmt.Println("novel name: ", novelName)
	tmpl, err := template.ParseFiles("www/search_result.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	result := searchNovel(novelName)
	err = tmpl.Execute(w, result)
	if err != nil {
		fmt.Println(err)
	}
	// go crawlNovel(novelName)

}
func pushHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		http.Redirect(w, req, "/404.html", 404)
		fmt.Println("bad requst")
		return
	}
	err := req.ParseForm()
	if err != nil {
		glog.Info(err)
	}

	email := req.FormValue("email")
	name := req.FormValue("novel")
	fmt.Fprint(w, `<h1>正在下载<< `, name, `>>,稍候将自动发送到您的kindle</h1>`)
	fmt.Println(email, name)
	go crawlNovel(name, email)
}
func main() {
	flag.Parse()
	http.Handle("/www/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", NotFoundHandler)
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/push", pushHandler)
	err := http.ListenAndServe("0.0.0.0:8081", nil)
	if err != nil {
		glog.Error(err)
	}
}

func download(url string) *string {
	resp, err := http.Get(url)
	if err != nil {
		glog.Errorln(err)
		time.Sleep(500 * time.Microsecond)
		return download(url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	str := ConvertGBKToUTF8(string(body))
	return &str
}
func ConvertGBKToUTF8(str string) string {
	dec := mahonia.NewDecoder("GBK")
	return dec.ConvertString(str)
}

func ConvertUTF8ToGBK(str string) string {
	enc := mahonia.NewEncoder("GBK")
	return enc.ConvertString(str)
}

func searchNovel(name string) []NovelInfo {
	realSearchURL := fmt.Sprintf(SEARCH_URL, ConvertUTF8ToGBK(name))
	uBody := download(realSearchURL)

	// parse href by goquery
	reader := strings.NewReader(*uBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		glog.Fatal(err)
		return nil
	}

	novelList := make([]NovelInfo, 0)
	doc.Find("#nr").Each(func(index int, s *goquery.Selection) {
		ni := NovelInfo{}
		s.Find("td").Each(func(td_index int, td_s *goquery.Selection) {
			switch td_index {
			case 0:
				ni.Name = td_s.Text()
				ni.EntryLink, _ = td_s.Find("a").Attr("href")
			case 1:
				ni.LatestChapter = td_s.Text()
			case 2:
				ni.Author = td_s.Text()
			}

		})

		novelList = append(novelList, ni)
	})
	return novelList
}
func crawlNovel(name, email string) {
	crawlerPath := filepath.Join(TOOLDIR, "crawl_novel")
	cmd := exec.Command(crawlerPath, "-name="+name)
	err := cmd.Run()
	if err != nil {
		glog.Error(err)
	}
	fmt.Println("下载完成: ", name)
	glog.Info("下载完成: ", name)

	ConvertTextToEpub(name)
	fmt.Println("转换epub完成")

	ConvertEpubToMobi(name)
	fmt.Println("转换mobi完成")

	SendMail(name, email)
	fmt.Println("发送邮件完成")
}
func ConvertEpubToMobi(name string) {
	tool := filepath.Join(TOOLDIR, "kindlegen")

	cmd := exec.Command(tool, name+".epub", "-o", name+".mobi")
	err := cmd.Run()
	if err != nil {
		glog.Error(err)
	}

}
func ConvertTextToEpub(name string) {
	novelPath := filepath.Join("novel", name)
	tool := filepath.Join(TOOLDIR, "create_epub.py")

	cmd := exec.Command("python", tool, "--dir="+novelPath)
	err := cmd.Run()
	if err != nil {
		glog.Error(err)
	}
}

func SendMail(name, email string) {
	mobiPath := name + ".mobi"
	tool := filepath.Join(TOOLDIR, "sendmail")

	cmd := exec.Command(tool, "-file", mobiPath, "-recipients", email)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		glog.Error(err)
	}
}
