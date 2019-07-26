package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type AutoGenerated struct {
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
			ContentType      string `json:"content-type"`
			TransferEncoding string `json:"transfer-encoding"`
			Date             string `json:"date"`
			Connection       string `json:"connection"`
		} `json:"res_header"`
		ResBody struct {
			Code    int    `json:"code"`
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    struct {
				ID          string      `json:"id"`
				Token       string      `json:"token"`
				Account     string      `json:"account"`
				Name        string      `json:"name"`
				Mobile      string      `json:"mobile"`
				CommAddress interface{} `json:"commAddress"`
				Email       interface{} `json:"email"`
				OrgID       interface{} `json:"orgId"`
				OrgName     interface{} `json:"orgName"`
				RoleIds     []string    `json:"roleIds"`
				RoleNames   []string    `json:"roleNames"`
				EntCodes    interface{} `json:"entCodes"`
				EntNames    interface{} `json:"entNames"`
				PushMsg     interface{} `json:"pushMsg"`
				Icon        interface{} `json:"icon"`
				ParkCode    int         `json:"parkCode"`
			} `json:"data"`
			Timestamp string `json:"timestamp"`
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
	url := os.Getenv("PLUGIN_URL")
	fmt.Printf("接口地址为: %s \n", url)
	date := getjson(url)
	generated := new(AutoGenerated)
	_ = json.Unmarshal([]byte(date), &generated)
	fmt.Printf("测试结果:  %s", generated.Message.Msg)
	if generated.Message.FailedNum != 0 {
		errormsg := "有错误接口共：" + string(generated.Message.FailedNum) + "个" + "。 请检查YAPI和接口"
		panic(errormsg)
	}
}

func getjson(url string) []byte {
	rep, err := http.Get(url)
	if err != nil || rep == nil {
		fmt.Printf("获取请求失败，错误信息： %s", err.Error())
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		fmt.Printf("解析json失败，错误信息： %s", err.Error())
	}
	return body
}
