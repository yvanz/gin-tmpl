# Gin Project Template

基于 [Gin](https://github.com/gin-gonic/gin) 框架以及 [XORM](https://xorm.io/) 构建的后端 API 服务。

项目结构规范依据 [Standard Go Project Layout](https://github.com/golang-standards/project-layout) 进行约束。



## 主要功能

- 项目整体使用 [golangci-lint](https://github.com/golangci/golangci-lint) 进行较为严格的静态代码检查，其配置文件为 `.golangci-lint.yml`
- `Gin` 实现的 WEB 服务
- 引入了跨域，记录请求和响应，pprof 以及 prometheus metrics 的中间件 
- API 参数校验及自定义翻译
- 基于 `XORM` 的 `MySQL` 读写分离
- 基于 [zap](https://github.com/uber-go/zap) 的日志处理
- 基于 [sarama](https://github.com/Shopify/sarama) 的 `Kafka` 生产者和消费者，可按需使用或删除。Kafka 集群最低支持 `0.8.2.0` 版本。
- swagger 文档



## 注意事项

### 关于配置文件
若没有配置

- `service_name`，则默认值为 gin-demo
- `local_ip`，则默认值为 0.0.0.0
- `api_port`，则默认值为 80

`run_mode` 为 debug，`Gin` 会运行在 debug 模式，在生产环境中，更换为任意值，则可以运行在 release 模式



### 关于编译

`vendor` 目录默认从项目中忽略，为加速 CI 中的编译速度，可以考虑将 vendor 添加到实际项目中



### 关于 Swagger 文档

当文档更新后，使用 `swag init -g ./cmd/app/main.go` 进行更新
