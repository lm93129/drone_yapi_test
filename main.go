package main

// PROJECT 项目id
// DATAURL 数据收集平台的地址
// post请求头会带上项目的id

import (
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

var YapiTestJson []YApiJSON

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
	// 发送数据到数据收集平台
	if os.Getenv("DATAURL") != "" {
		err := gout.
			POST(os.Getenv("DATAURL")).
			SetHeader(gout.H{"project_id": os.Getenv("PROJECT")}).
			SetJSON(YapiTestJson).
			Do()

		if err != nil {
			log.Printf("发送失败：%s\n", err)
		}
	}

}

func CheckApi(BaseUrl, id string) {
	l := strings.Split(id, ",")
	c := make(chan int, len(l))
	v := 0

	// 启动多线程遍历请求每一个用例集合
	for _, v := range l {
		url := BaseUrl + v
		go func(u, v string) {
			YapiAutoTest(u, v, c)
		}(url, v)
	}

	// 打印每一个用例集合
	for i := 0; i < len(l); i++ {
		tm := time.NewTimer(time.Second * 20)
		select {
		case msg := <-c:
			v += msg
		case <-tm.C:
			log.Println("测试超时，请检查网络环境")
		}
	}

	// 做个统计，在所有接口测试完毕之后统计是否有不通过的测试用例集合
	// 如果有不通过的用例就会报错跳出
	if v > 0 {
		log.Panic("接口测试不通过，请检查")
	}

}

func YapiAutoTest(url, v string, c chan int) {
	apiJson := YApiJSON{}
	err := gout.GET(url).
		SetTimeout(20 * time.Second).
		BindJSON(&apiJson).
		Do()

	if err != nil {
		log.Printf("Yapi请求错误: %s\n", err)
	}

	log.Printf("开始测试ID为 %s 的用例集合 ", v)

	log.Println(apiJson.Message.Msg)

	// 将个测试用例集合的测试结果收集
	YapiTestJson = append(YapiTestJson, apiJson)

	for i := 0; i < len(apiJson.List); i++ {

		name := apiJson.List[i].Name
		message := apiJson.List[i].ValidRes[0].Message
		log.Printf("接口用例名称：%s , 验证结果： %s \n", name, message)
	}

	if apiJson.Message.FailedNum != 0 {
		log.Printf("接口验证不通过，错误数：%d , 用例集合耗时：%s \n", apiJson.Message.FailedNum, apiJson.RunTime)
		c <- 1
	}
	c <- 0

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
		ResHeader interface{} `json:"res_header"`
		ResBody   interface{} `json:"res_body"`
		Params    interface{} `json:"params"`
	} `json:"list"`
}
