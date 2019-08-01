package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
)

func main() {
	url := os.Getenv("PLUGIN_URL")
	re, _ := regexp.Compile(`&mode=json`)
	match := re.FindAllStringSubmatch(url, -1)
	if match != nil {
		panic("请使用html的接口")
	}
	fmt.Printf("YAPI地址为: %s \n", url)
	List(url)
}

func List(url string) {
	var htmldate = html(url)
	for _, m := range statistical(htmldate) {
		fmt.Printf("%s  %s%s\n", m[1], m[2], m[3])
	}
	for _, m := range getlist(htmldate) {
		fmt.Printf("%s   %s\n", m[1], m[2])
	}
	for _, stat := range statistical(htmldate) {
		i := string(stat[2])
		if i != "" {
			panic("接口测试不通过，请检查YAPI和接口")
		}
	}
}

func html(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	s, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Panic(err)
	}
	return s
}

func getlist(htmldate []byte) [][][]byte {
	re := regexp.MustCompile(`href="#[0-9]+">([^<]+)</a>
    <div title="([^\"]+)"`)
	match := re.FindAllSubmatch(htmldate, -1)
	return match
}

func statistical(htmldate []byte) [][][]byte {
	re := regexp.MustCompile(`<div class="summary"><div>([^<]+)<span class="success"> ([0-9]+)</span>([^<]+)</div>`)
	match := re.FindAllSubmatch(htmldate, -1)
	return match
}
