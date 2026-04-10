# Gin Backend Skill 文档

## 项目概述
基于 Gin + GORM 的后端微服务脚手架，采用分层架构。

## 目录结构
```
cmd/app/main.go           # 入口
internal/
  app/run.go             # 启动配置 (cobra)
  handler/               # Controller 层
  logic/                 # Service 层
  common/                # 公共组件 (BaseController, 错误码)
  config/                # 配置加载
models/
  dbDemo.go              # 数据模型 (含 Meta 公共字段)
  repo/demo.go           # Repo 接口定义
  factory/demo.go        # Repo 实现 (GORM CRUD)
pkg/
  apiserver/             # HTTP Server 封装
  gormdb/                # 数据库封装
  logger/                # 日志
  middleware/            # 中间件
  kafka/                 # Kafka
  rediscache/            # Redis
  tracer/                # Jaeger
```

## 新增业务模块步骤

### 1. 定义 Model
文件: `models/dbXXX.go`
```go
type XXX struct {
    Meta
    FieldName string `json:"field_name" gorm:"column:field_name"`
}
func (m *XXX) ColumnFieldName() string { return "field_name" }
```

### 2. 定义 Repo 接口
文件: `models/repo/xxx.go`
```go
type XXXRepo interface {
    gormdb.GetListCrud
    gormdb.GetByIDCrud
    gormdb.CreateCrud
    gormdb.UpdateCrud
    Deletes([]int64) error
}
```

### 3. 实现 Repo
文件: `models/factory/xxx.go`
```go
type xxxCrudImpl struct { Conn *gorm.DB }

func XXXRepo(db *gorm.DB) repo.XXXRepo {
    return &xxxCrudImpl{Conn: db}
}

func (r *xxxCrudImpl) GetList(q gormdb.BasicQuery, model, list interface{}) (int64, error) {
    crud := gormdb.NewCRUD(r.Conn)
    return crud.GetList(q, model, list)
}
// ... 实现其他方法
```

### 4. 定义 Logic Service
文件: `internal/logic/srvxxx/srv_xxx.go`
```go
type Svc struct {
    Ctx context.Context
    ID  int64
}

func (s *Svc) getRepo() repo.XXXRepo {
    return factory.XXXRepo(gormdb.Cli(s.Ctx))
}

// 添加业务方法
func (s *Svc) GetList(q gormdb.BasicQuery) (*common.ListData, error) {
    data := &common.ListData{PageOffset: q.Offset, PageLimit: q.Limit}
    crud := s.getRepo()
    list := make([]models.XXX, 0)
    total, err := crud.GetList(q, &models.XXX{}, &list)
    if err != nil { return nil, common.NewCodeWithErr(common.ErrorDatabaseRead, err) }
    data.Counts = total
    data.Data = list
    return data, nil
}
```

### 5. 定义 Handler
文件: `internal/handler/xxx.go`
```go
type xxxController struct { common.BaseController }

func newXXXController(base common.BaseController) *xxxController {
    return &xxxController{BaseController: base}
}

// GET 列表
func (c *xxxController) Get(ctx *gin.Context) {
    var svc srvxxx.Svc
    q := gormdb.BasicQuery{Offset: ..., Limit: ..., Query: ...}
    data, err := svc.GetList(q)
    c.Response(ctx, data, err)
}

// POST 创建
func (c *xxxController) Create(ctx *gin.Context) {
    var params AddParams  // 定义在 logic 中
    if !c.CheckParams(ctx, &params) { return }
    var svc srvxxx.Svc
    err := svc.Add(params)
    c.Response(ctx, nil, err)
}

// PUT 更新
func (c *xxxController) Update(ctx *gin.Context) {
    var params AddParams
    if !c.CheckParams(ctx, &params) { return }
    var svc srvxxx.Svc
    if svc.ID, isNum = c.CheckNumber(ctx, ctx.Param("id")); !isNum { return }
    err := svc.Mod(params)
    c.Response(ctx, nil, err)
}

// DELETE 删除
func (c *xxxController) Delete(ctx *gin.Context) {
    ids := strings.Split(ctx.Param("ids"), ",")
    err := srvxxx.Svc{}.Delete(ids)
    c.Response(ctx, nil, err)
}
```

### 6. 注册路由
文件: `internal/handler/routers.go`
```go
func RegisterHandler(...) {
    v1API := engine.Group("/api/v1")
    group := v1API.Group("/xxx")
    ctrl := newXXXController(base)
    group.GET("", ctrl.Get)
    group.GET("/:id", ctrl.GetByID)
    group.POST("", ctrl.Create)
    group.PUT("/:id", ctrl.Update)
    group.DELETE("/:ids", ctrl.Delete)
}
```

## 代码规范
- Model 字段使用 `json` + `gorm:"column:xxx"` 双标签
- 使用 validator 做参数校验 (binding:"required")
- 错误处理: `common.NewCodeWithErr(code, err)`
- 响应格式: `{ret_code, data_set, message}`
- 不添加无意义注释，代码自解释