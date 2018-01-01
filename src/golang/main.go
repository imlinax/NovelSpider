package main

import (
	"path/filepath"
	"os"
	"strings"
	"fmt"
	"github.com/golang/glog"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"io/ioutil"
)
var (
	index = 1
	dir = "unkown"
)
func main() {
	novel_name := "圣墟"
	dir = novel_name
	os.Mkdir(dir, 0755)

	start_url :="http://www.biqiuge.com/book/4772/2940354.html"
	next_url := &start_url
	for {
		next_url = downloadPage(next_url)
		if next_url == nil{
			break
		}
	}
}
func downloadPage(url *string) (next *string){
	doc,err := goquery.NewDocument(*url)
	if err != nil {
		glog.Error(err)
		return nil
	}
	dec := mahonia.NewDecoder("GB18030")
	title := doc.Find(".bookname").Find("h1").Text()
	title = dec.ConvertString(title)
	fmt.Println(title)
	doc.Find("#content").Each(func(i int,s *goquery.Selection){
		str := dec.ConvertString(s.Text())
		str = strings.Replace(str, "聽"," ",-1)
		nt := fmt.Sprintf("%4d %s.txt", index, title)
		path := filepath.Join(dir,nt)
		index++
		ioutil.WriteFile(path, []byte(str), 0644)
		//fmt.Println(str)
	})
	doc.Find(".bottem").Each(func(i int,s *goquery.Selection){
		s.Find("a").Each(func(i int, c *goquery.Selection){
			text := dec.ConvertString(c.Text())
			if text == "下一章" {
				href,ok := c.Attr("href")
				if ok {
					fmt.Println("下一章： ", href)
					next = &href
					return
				}
			}
		})
	})
	return
}