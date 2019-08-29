请求yapi的接口，进行测试集合的测试。使用方法如下：
```yaml
 #YAPI接口自动测试
  - name: yapi-test
    image: registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
    settings:
      host: http://yapi.com.cn/
      id: [11, 31, 15]
      token: ab1058076e0945cf14
      env: env_11=网关地址
      mode: json
```
使用docker运行
```bash
docker run --rm \
  -e PLUGIN_HOST=http://yapi.com.cn/ \
  -e PLUGIN_MODE=json \
  -e PLUGIN_TOKEN=ab1058076e0945cf14 \
  -e PLUGIN_ID=[11, 31, 15] \
  -e PLUGIN_ENV="env_11=网关地址" \
  -v $(pwd):/drone/src \
  -w /drone/src \
  registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
```
注意：host地址后面必须接上/