package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/golang/glog"
)

var (
	dir       = "unkown"
	novelName = flag.String("name", "", "novel name")
)

const (
	SEARCH_SITE = `http://www.biquge.com.tw`
	SEARCH_URL  = SEARCH_SITE + `/modules/article/soshu.php?searchkey=%s`
)

func getNovelIndex(name string) {
	realSearchURL := fmt.Sprintf(SEARCH_URL, ConvertUTF8ToGBK(name))
	resp, err := http.Get(realSearchURL)
	if err != nil {
		glog.Errorln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorln(err)
	}
	uBody := ConvertGBKToUTF8(string(body))
	//fmt.Println(uBody)

	// parse href by goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(uBody))
	if err != nil {
		glog.Fatal(err)
	}
	var wg sync.WaitGroup
	doc.Find("#list").Find("dd").Each(func(index int, s *goquery.Selection) {
		chapter := s.Children()
		fmt.Print(chapter.Text())
		href, ok := chapter.Attr("href")
		if ok {
			url := SEARCH_SITE + href
			fmt.Println(url)
			wg.Add(1)
			go func() {
				defer wg.Done()
				downloadPage(&url, index)
			}()
		} else {
			fmt.Println("no link find")
		}
	})
	wg.Wait()
}
func main() {
	flag.Parse()
	if len(*novelName) == 0 {
		flag.Usage()
		os.Exit(0)
	}
	os.Mkdir("novel", 0755)
	dir = filepath.Join("novel", *novelName)
	os.Mkdir(dir, 0755)
	getNovelIndex(*novelName)
}
func downloadPage(url *string, index int) (next *string) {
	doc, err := goquery.NewDocument(*url)
	if err != nil {
		glog.Error(err)
		time.Sleep(time.Millisecond * 500)
		return downloadPage(url, index)
	}
	dec := mahonia.NewDecoder("GB18030")
	title := doc.Find(".bookname").Find("h1").Text()
	title = dec.ConvertString(title)
	doc.Find("#content").Each(func(i int, s *goquery.Selection) {
		str := dec.ConvertString(s.Text())
		str = strings.Replace(str, "聽", " ", -1)
		nt := fmt.Sprintf("%4d %s.txt", index, title)
		fmt.Println(nt)
		path := filepath.Join(dir, nt)
		ioutil.WriteFile(path, []byte(str), 0644)
	})
	return
}
