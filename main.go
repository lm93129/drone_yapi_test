package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	dir, _ := os.Getwd()
	_ = godotenv.Load(dir + "/.env")

	YApiHOST := os.Getenv("PLUGIN_HOST")
	token := os.Getenv("PLUGIN_TOKEN")
	id := os.Getenv("PLUGIN_ID")
	env := os.Getenv("PLUGIN_ENV")
	BaseURL := fmt.Sprintf("%s/api/open/run_auto_test?token=%s&%s&mode=json&email=false&download=false&id=",
		YApiHOST,
		token,
		env,
	)

	CheckApi(BaseURL, id)
}

func CheckApi(BaseUrl, id string) {
	l := strings.Split(id, ",")
	c := make(chan int, len(l))
	v := 0
	for _, v := range l {
		url := BaseUrl + v
		go func(u, v string) {
			YapiAutoTest(u, v, c)
		}(url, v)
	}
	for i := 0; i < len(l); i++ {
		tm := time.NewTimer(time.Second * 20)
		select {
		case msg := <-c:
			v += msg
		case <-tm.C:
			log.Println("测试超时，请检查网络环境")
		}
	}
	if v > 0 {
		log.Panic("接口测试不通过，请检查")
	}
}

var myClient = &http.Client{Timeout: 20 * time.Second}

func YapiAutoTest(url, v string, c chan int) {
	apiJSON := new(YApiJSON)
	r, err := myClient.Get(url)
	dropErr(err)
	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(apiJSON)
	log.Printf("开始测试ID为 %s 的用例集合 ", v)
	log.Println(apiJSON.Message.Msg)

	for i := 0; i < len(apiJSON.List); i++ {
		name := apiJSON.List[i].Name
		message := apiJSON.List[i].ValidRes[0].Message
		log.Printf("接口用例名称：%s , 验证结果： %s \n", name, message)
	}
	if apiJSON.Message.FailedNum != 0 {
		log.Printf("接口验证不通过，错误数：%d , 耗时：%s \n", apiJSON.Message.FailedNum, apiJSON.RunTime)
		c <- 1
	}
	c <- 0
}

func dropErr(e error) {
	if e != nil {
		log.Panic(e)
	}
}

// YApiJSON 返回的json序列化
type YApiJSON struct {
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
