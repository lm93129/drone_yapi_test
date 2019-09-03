请求yapi的接口，进行测试集合的测试。使用方法如下：
```yaml
 #YAPI接口自动测试
  - name: yapi-test
    image: registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
    settings:
      host: http://yapi.com.cn/
      id: [11,31,15]
      token: ab1058076e0945cf14
      env: env_11=网关地址
```
使用docker运行
```bash
docker run --rm \
  -e PLUGIN_HOST=http://yapi.com.cn \
  -e PLUGIN_TOKEN=ab1058076e0945cf14 \
  -e PLUGIN_ID="11,31,15" \
  -e PLUGIN_ENV="env_11=网关地址" \
  registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
```

使用gitlab-ci运行
```yaml
apitest:
  image: registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
  stage: apitest
  variables:
    PLUGIN_HOST: "http://yapi.com.cn"
    PLUGIN_TOKEN: "ab1058076e0945cf14"
    PLUGIN_ID: "11,31,15"
    PLUGIN_ENV: "env_11=网关地址"
  script:
    - /bin/apitest
```