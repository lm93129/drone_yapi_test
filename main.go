package main

// PLUGIN_PROJECT 项目id
// PLUGIN_DATAURL 数据收集平台的地址
// post请求头会带上项目的project_id：PLUGIN_PROJECT

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/guonaihong/gout"
	"github.com/joho/godotenv"
)

var YapiTestJson []YApiJSON

func main() {
	// 获取当前目录
	dir, _ := os.Getwd()
	// 自动载入环境变量
	_ = godotenv.Load(dir + "/.env")
	// yapi地址
	YApiHOST := os.Getenv("PLUGIN_HOST")
	// 令牌
	token := os.Getenv("PLUGIN_TOKEN")
	// 测试ID
	id := os.Getenv("PLUGIN_ID")
	// 环境变量选项
	env := os.Getenv("PLUGIN_ENV")
	// 子集选项
	descendants := os.Getenv("PLUGIN_DESCENDANTS")
	BaseURL := fmt.Sprintf("%s/api/open/run_auto_test?token=%s&%s&mode=json&email=false&download=false&descendants=%s&id=",
		YApiHOST,
		token,
		env,
		descendants,
	)

	check := CheckApi(BaseURL, id)

	// 如果存在PLUGIN_DATAURL则发送数据到数据收集平台
	if os.Getenv("PLUGIN_DATAURL") != "" {
		err := gout.
			POST(os.Getenv("PLUGIN_DATAURL")).
			SetHeader(gout.H{"project_id": os.Getenv("PLUGIN_PROJECT")}).
			SetJSON(YapiTestJson).
			Do()

		if err != nil {
			log.Printf("发送失败：%s\n", err)
		} else {
			log.Println("数据发送成功")
		}
	}

	if check > 0 {
		log.Panicf("接口测试不通过，请检查接口。错误用例数共：%d 个", check)
	}
}

func CheckApi(BaseUrl, id string) int {
	l := strings.Split(id, ",")
	c := make(chan int, len(l))
	var errCase int

	// 启动多线程遍历请求每一个用例集合
	for _, v := range l {
		url := BaseUrl + v
		go func(u, v string) {
			errCase = errCase + YapiAutoTest(u, v)
			c <- 1
		}(url, v)
	}

	// 打印每一个用例集合
	tm := time.NewTimer(time.Second * 20)
	for range l {
		select {
		case <-c:
		case <-tm.C:
			log.Println("测试超时，请检查网络环境")
		}
	}

	// 返回总的错误用例数
	return errCase
}

func YapiAutoTest(url, v string) int {
	apiJson := YApiJSON{}
	// 请求Yapi的测试
	err := gout.GET(url).
		SetTimeout(20 * time.Second).
		BindJSON(&apiJson).
		Do()

	// 如果错误，则打印错误出来
	if err != nil {
		log.Printf("Yapi请求错误: %s\n", err)
	}

	log.Printf("开始测试ID为 %s 的用例集合 ", v)

	log.Printf("%s 用例集合耗时：%s \n", apiJson.Message.Msg, apiJson.RunTime)

	// 将个测试用例集合的测试结果收集
	YapiTestJson = append(YapiTestJson, apiJson)

	for i := 0; i < len(apiJson.List); i++ {
		name := apiJson.List[i].Name
		message := apiJson.List[i].ValidRes[0].Message
		if message != "验证通过" {
			log.Printf("接口用例名称：%s , 验证结果： %s \n", name, message)
		}
	}

	if apiJson.Message.FailedNum != 0 {
		log.Printf("接口验证不通过，错误数：%d , 用例集合耗时：%s \n", apiJson.Message.FailedNum, apiJson.RunTime)
	}

	return apiJson.Message.FailedNum
}

// YApi 返回的json序列化结构体
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
		ResHeader interface{} `json:"res_header"`
		ResBody   interface{} `json:"res_body"`
		Params    interface{} `json:"params"`
	} `json:"list"`
}
