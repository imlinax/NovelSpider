package main

import (
	"github.com/axgle/mahonia"
)

func ConvertGBKToUTF8(str string) string {
	dec := mahonia.NewDecoder("GBK")
	return dec.ConvertString(str)
}

func ConvertUTF8ToGBK(str string) string {
	enc := mahonia.NewEncoder("GBK")
	return enc.ConvertString(str)
}
