请求 `yapi` 的接口，进行测试集合的测试。使用方法如下：

如果项目目录中有`.env`配置文件，则下面的命令不需要加其他参数，会自动加载项目中.env中的配置，方便每个项目自行维护
`.env`文件的配置见`.env.example`中的示例

## 需要注意的地方：
> 测试用例集合的测试时间超过20秒会出现超时
> 测试用例集合出现错误，程序不会直接退出，会测试完所有的用例集合后再检测是否有错误的集合，然后再退出
> YAPI上的测试用例中有test的断言检测，如果检测失败，会报错，并且会打印错误到日志

```yaml
#YAPI接口自动测试
- name: yapi-test
  image: registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
  settings:
    host: http://yapi.com.cn/
    id: [11,31,15]
    token: ab1058076e0945cf14
    env: env_11=网关地址
    # 是否包含子集，适用于yapi的子集分支,如果没有该参数，则默认不包含子集
    DESCENDANTS: "true"
    # 非必须参数，该参数会发送测试过程中的yapi参数到DATAURL地址中，进行数据收集
    DATAURL: http://127.0.0.1:3000/interface/apidata
    PROJECT: etl-root
```

使用 docker 运行，如果是用env文件运行的，需要挂载env文件到插件目录并使用-w指定env所在目录

```bash
docker run --rm \
  -w /src/xxx项目(仅当使用env文件的时候需要) \
  -v /src/xxx项目(仅当使用env文件的时候需要) \
  -e PLUGIN_HOST=http://yapi.com.cn \
  -e PLUGIN_TOKEN=ab1058076e0945cf14 \
  -e PLUGIN_ID="11,31,15" \
  -e PLUGIN_ENV="env_11=网关地址" \
  # 是否包含子集，适用于yapi的子集分支
  -e PLUGIN_DESCENDANTS="true" \
  # 非必须参数，该参数会发送测试过程中的yapi参数到DATAURL地址中，进行数据收集
  # -e PLUGIN_DATAURL="http://127.0.0.1:3000/interface/apidata" \
  # -e PLUGIN_PROJECT="etl-root" \
  registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
```

使用 gitlab-ci 运行

```yaml
apitest:
  image: registry.cn-hangzhou.aliyuncs.com/lm93129/drone_yapi_test:latest
  stage: apitest
  variables:
    PLUGIN_HOST: "http://yapi.com.cn"
    PLUGIN_TOKEN: "ab1058076e0945cf14"
    PLUGIN_ID: "11,31,15"
    PLUGIN_ENV: "env_11=网关地址"
    # 是否包含子集，适用于yapi的子集分支
    PLUGIN_DESCENDANTS: "true"
    # 非必须参数，该参数会发送测试过程中的yapi参数到DATAURL地址中，进行数据收集
    PLUGIN_DATAURL: "http://127.0.0.1:3000/interface/apidata"
    PLUGIN_PROJECT: "etl-root"
  script:
    - /bin/apitest
```
发送yapi测试结果功能需要自己有个数据收集的平台，改功能会上传yapi的测试结果json数据到你指定的地址中，PROJECT项目名称会在header的project_id字段中，所以自建的测试数据收集平台需要获取发送过来的json，header中的project_id字段内容。
具体发送的json示例，可以在wiki中看到