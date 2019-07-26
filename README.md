请求yapi的接口，进行测试集合的测试。使用方法如下：
```
 #YAPI接口自动测试
  - name: yapi-test
    image: registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
    settings:
      url: http://yapi.com/api/open/run_auto_test?id=11&token=ab1058e0945cf14&env_11=测试环境&mode=json&email=false&download=false
```

注意yapi的测试地址接口必须选择json和非下载