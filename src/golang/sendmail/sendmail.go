package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/axgle/mahonia"
	"github.com/golang/glog"
	"gopkg.in/gomail.v2"
)

var (
	flagNovelFile  = flag.String("file", "", "novel(*.mobi) file path")
	flagRecipients = flag.String("recipients", "", "recipients")
)

func main() {
	flag.Parse()

	if len(*flagNovelFile) == 0 || len(*flagRecipients) == 0 {
		flag.Usage()
		os.Exit(-1)
	}
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "wangli.qufu@163.com", "王立")
	msg.SetHeader("To", msg.FormatAddress(*flagRecipients, *flagRecipients))

	msg.SetHeader("Subject", filepath.Base(*flagNovelFile))
	msg.SetBody("text/html", "Hello <h1> let me see you </h1>")

	h := make(map[string][]string, 0)
	//h["Content-Type"] = []string{`application/octet-stream`} //要设置这个，否则中文会乱码
	filename := mahonia.NewEncoder("gbk").ConvertString(*flagNovelFile)
	h["Content-Disposition"] = []string{`attachment; charset=utf-8;filename=` + filename} //要设置这个，否则中文会乱码
	fileSetting := gomail.SetHeader(h)

	msg.Attach(*flagNovelFile, fileSetting)

	mailer := gomail.NewDialer("smtp.163.com", 25, "wangli.qufu@163.com", "**********")
	if err := mailer.DialAndSend(msg); err != nil {
		glog.Fatal(err)
	}

}
