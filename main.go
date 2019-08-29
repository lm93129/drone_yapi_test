package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"time"
)

type YApiJson struct {
	Message struct {
		Msg        string `json:"msg"`
		Len        int    `json:"len"`
		SuccessNum int    `json:"successNum"`
		FailedNum  int    `json:"failedNum"`
	} `json:"message"`
	RunTime string `json:"runTime"`
	Numbs   int    `json:"numbs"`
	List    []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		Code     int    `json:"code"`
		ValidRes []struct {
			Message string `json:"message"`
		} `json:"validRes"`
		Status  int    `json:"status"`
		URL     string `json:"url"`
		Method  string `json:"method"`
		Headers struct {
			ContentType string `json:"Content-Type"`
			Accept      string `json:"accept"`
		} `json:"headers"`
		ResHeader struct {
			Server           string `json:"server"`
			Date             string `json:"date"`
			ContentType      string `json:"content-type"`
			TransferEncoding string `json:"transfer-encoding"`
			Connection       string `json:"connection"`
			Vary             string `json:"vary"`
		} `json:"res_header"`
		ResBody struct {
			Code      int         `json:"code"`
			Status    bool        `json:"status"`
			Message   string      `json:"message"`
			Data      interface{} `json:"data"`
			Timestamp string      `json:"timestamp"`
		} `json:"res_body"`
		Params struct {
			Username    string `json:"username"`
			Password    string `json:"password"`
			DeviceType  string `json:"deviceType"`
			ContentType string `json:"Content-Type"`
			Accept      string `json:"accept"`
		} `json:"params"`
	} `json:"list"`
}

func main() {
	YApihost := os.Getenv("PLUGIN_HOST")
	mode := os.Getenv("PLUGIN_MODE")
	token := os.Getenv("PLUGIN_TOKEN")
	id := os.Getenv("PLUGIN_ID")
	env := os.Getenv("PLUGIN_ENV")
	beseurl := YApihost +
		"api/open/run_auto_test?" +
		mode + "&" +
		token + "&" +
		env + "&email=false&download=false&"

	checkapi(beseurl, id, mode)
}

func checkapi(baseurl string, id string, mode string) {
	i := 0
	for _, v := range strings.Split(id, ",") {
		url := baseurl + v
		if mode == "json" {
			fmt.Printf("测试接口地址为: %s \n", url)
			i = i + judge_json(url)
		} else {
			fmt.Printf("测试接口地址为: %s \n", url)
			i = i + testapi_html(url)
		}
		if i > 0 {
			log.Panic("接口测试不通过，请检查YApi和接口")
		}
	}
}

func testapi_html(url string) int {
	var htmldate = getbody(url)
	for _, m := range statistical(htmldate) {
		fmt.Printf("%s  %s%s\n", m[1], m[2], m[3])
	}
	for _, m := range getlist(htmldate) {
		fmt.Printf("%s   %s\n", m[1], m[2])
	}
	for _, stat := range statistical(htmldate) {
		i := string(stat[2])
		if i != "" {
			return 1
		}
	}
	return 0
}

func getbody(url string) []byte {
	resp, err := http.Get(url)
	dropErr(err)
	defer resp.Body.Close()
	s, err := httputil.DumpResponse(resp, true)
	dropErr(err)
	return s
}

func getlist(contents []byte) [][][]byte {
	re := regexp.MustCompile(`href="#[0-9]+">([^<]+)</a>
    <div title="([^\"]+)"`)
	match := re.FindAllSubmatch(contents, -1)
	return match
}

func statistical(contents []byte) [][][]byte {
	re := regexp.MustCompile(`<div class="summary"><div>([^<]+)<span class="success"> ([0-9]+)</span>([^<]+)</div>`)
	match := re.FindAllSubmatch(contents, -1)
	return match
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	dropErr(err)
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func judge_json(url string) int {
	apijs := new(YApiJson)
	err := getJson(url, apijs)
	if err != nil {
		log.Panicln("json解析失败")
	}
	fmt.Println(apijs.Message.Msg)
	for i := 0; i < len(apijs.List); i++ {
		name := apijs.List[i].Name
		message := apijs.List[i].ValidRes[0].Message
		fmt.Printf("接口用例名称：%s , 验证结果： %s \n", name, message)
	}
	if apijs.Message.FailedNum != 0 {
		fmt.Printf("接口验证不通过，错误数：%d , 耗时：%s \n", apijs.Message.FailedNum, apijs.RunTime)
		return 1
	}
	return 0
}

func dropErr(e error) {
	if e != nil {
		log.Panic(e)
	}
}
