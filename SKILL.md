---
name: gin-tmpl
description: 基于 Gin + GORM 的 Go 后端脚手架，使用项目源码作为模板生成 API
---

# Gin Backend Skill 文档

## 使用说明

代码模板为项目本身的源码，使用时参考以下文件：
- **Model**: `models/dbDemo.go`
- **Repo 接口**: `models/repo/demo.go`
- **Repo 实现**: `models/factory/demo.go`
- **Logic Service**: `internal/logic/srvdemo/srv_demo.go`
- **Handler**: `internal/handler/demo.go`
- **路由注册**: `internal/handler/routers.go`

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
  dbXXX.go              # 数据模型 (含 Meta 公共字段)
  repo/xxx.go           # Repo 接口定义
  factory/xxx.go        # Repo 实现 (GORM CRUD)
pkg/
  apiserver/             # HTTP Server 封装
  gormdb/                # 数据库封装
  logger/                # 日志
  middleware/            # 中间件
  kafka/                 # Kafka
  rediscache/            # Redis
  tracer/                # Jaeger
```

## GORM 泛型接口
项目已使用 Go 1.25 泛型，以下接口可直接嵌入使用：
```go
// gormdb/basic.go
type GetListCrud[T any] interface {
    GetList(q BasicQuery, model *T, list *[]T) (total int64, err error)
}
type GetByIDCrud[T any] interface {
    GetByID(model *T, id int64) error
}
type GetOneByConCrud[T any] interface {
    GetOneByCon(con any, model *T, args ...any) error
}
type FindByConCrud[T any] interface {
    FindByCon(con any, model *T, args ...any) error
}
type CreateCrud[T any] interface {
    Create(model *T) error
}
type UpdateCrud[T any] interface {
    UpdateWithMap(model *T, u map[string]any) error
}
type DeleteCrud[T any] interface {
    Delete(model *T, hardDelete bool) error
}
type BasicCrud[T any] interface {
    GetListCrud[T]
    GetByIDCrud[T]
    GetOneByConCrud[T]
    FindByConCrud[T]
    CreateCrud[T]
    UpdateCrud[T]
    DeleteCrud[T]
}
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

**⚠️ 定义 Model 后必须注册到 AllTables**，在 `models/init.go` 中：
```go
var AllTables = []any{
    &Demo{},
    &XXX{}, // 新增
}
```

AllTables 在服务启动时会通过 `apiserver.Migration(models.AllTables)` 自动执行 GORM AutoMigrate。若不注册，新表/新字段不会被迁移到数据库。

### 2. 定义 Repo 接口
文件: `models/repo/xxx.go`
```go
type XXXRepo[T any] interface {
    gormdb.GetListCrud[T]
    gormdb.GetByIDCrud[T]
    gormdb.CreateCrud[T]
    gormdb.UpdateCrud[T]
    Deletes([]int64) (err error)
}
```

### 3. 实现 Repo
文件: `models/factory/xxx.go`
```go
type xxxCrudImpl[T any] struct {
    Conn *gorm.DB
}

func XXXRepo[T any](db *gorm.DB) repo.XXXRepo[T] {
    return &xxxCrudImpl[T]{Conn: db}
}

func (r *xxxCrudImpl[T]) GetList(q gormdb.BasicQuery, model *T, list *[]T) (total int64, err error) {
    crud := gormdb.NewCRUD[T](r.Conn)
    return crud.GetList(q, model, list)
}
// ... 实现其他方法
```

### 4. 定义 Logic Service
文件: `internal/logic/srvxxx/srv_xxx.go`

包导入: `github.com/yvanz/gin-tmpl/internal/logic/srvxxx`

```go
package srvxxx

import (
    "github.com/yvanz/gin-tmpl/internal/common"
    "github.com/yvanz/gin-tmpl/models"
    "github.com/yvanz/gin-tmpl/models/factory"
    "github.com/yvanz/gin-tmpl/models/repo"
    "github.com/yvanz/gin-tmpl/pkg/gormdb"
)

type Svc struct {
    Ctx context.Context
    ID  int64
}

func (s *Svc) getRepo() repo.XXXRepo[models.XXX] {
    return factory.XXXRepo[models.XXX](gormdb.Cli(s.Ctx))
}

// 请求参数定义
type AddParams struct {
    UserName string `json:"user_name" binding:"required"`
}

// GET 列表
func (s *Svc) GetList(q gormdb.BasicQuery) (*common.ListData, error) {
    data := &common.ListData{PageOffset: q.Offset, PageLimit: q.Limit}
    crud := s.getRepo()
    var table models.XXX
    list := make([]models.XXX, 0)
    total, err := crud.GetList(q, &table, &list)
    if err != nil { return nil, common.NewCodeWithErr(common.ErrorDatabaseRead, err) }
    data.Counts = total
    data.Data = list
    return data, nil
}

// GET 单条
func (s *Svc) GetByID() (*models.XXX, error) {
    crud := s.getRepo()
    d := &models.XXX{}
    err := crud.GetByID(d, s.ID)
    if err != nil { return nil, common.NewCodeWithErr(common.ErrorDatabaseRead, err) }
    return d, nil
}

// POST 新增
func (s *Svc) Add(params AddParams) error {
    crud := s.getRepo()
    d := &models.XXX{UserName: params.UserName}
    if err := crud.Create(d); err != nil {
        return common.NewCodeWithErr(common.ErrorDatabaseWrite, err)
    }
    return nil
}

// PUT 更新
func (s *Svc) Mod(params AddParams) error {
    crud := s.getRepo()
    d := &models.XXX{}
    if err := crud.GetByID(d, s.ID); err != nil {
        return common.NewCodeWithErr(common.ErrorDatabaseRead, err)
    }
    u := map[string]any{d.ColumnUserName(): params.UserName}
    if err := crud.UpdateWithMap(d, u); err != nil {
        return common.NewCodeWithErr(common.ErrorDatabaseWrite, err)
    }
    return nil
}

// DELETE 删除
func (s *Svc) Delete(ids []string) error {
    // id 转换逻辑...
    return crud.Deletes(idList)
}
```

### 5. 定义 Handler
文件: `internal/handler/xxx.go`

包导入: `github.com/yvanz/gin-tmpl/internal/handler` + `github.com/yvanz/gin-tmpl/internal/logic/srvxxx`

```go
type xxxController struct { common.BaseController }


func newXXXController(base common.BaseController) *xxxController {
    return &xxxController{BaseController: base}
}

// GET 列表
func (c *xxxController) Get(ctx *gin.Context) {
    var svc srvxxx.Svc
    q := gormdb.BasicQuery{Offset: ..., Limit: ..., Query: ...}
    svc.Ctx = ctx  // 传入 Gin context
    data, err := svc.GetList(q)
    c.Response(ctx, data, err)
}

// GET 单条
func (c *xxxController) GetByID(ctx *gin.Context) {
    var svc srvxxx.Svc
    if svc.ID, isNum = c.CheckNumber(ctx, ctx.Param("id")); !isNum { return }
    svc.Ctx = ctx
    data, err := svc.GetByID()
    c.Response(ctx, data, err)
}

// POST 创建
func (c *xxxController) Create(ctx *gin.Context) {
    var params srvxxx.AddParams
    if !c.CheckParams(ctx, &params) { return }
    var svc srvxxx.Svc
    svc.Ctx = ctx
    err := svc.Add(params)
    c.Response(ctx, nil, err)
}

// PUT 更新
func (c *xxxController) Update(ctx *gin.Context) {
    var params srvxxx.AddParams
    if !c.CheckParams(ctx, &params) { return }
    var svc srvxxx.Svc
    if svc.ID, isNum = c.CheckNumber(ctx, ctx.Param("id")); !isNum { return }
    svc.Ctx = ctx
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


函数签名: `func RegisterHandler(tra opentracing.Tracer, engine *gin.Engine)`


```go
import (
    "github.com/yvanz/gin-tmpl/internal/common"
    "github.com/yvanz/gin-tmpl/pkg/middleware"
)

var base common.BaseController

func RegisterHandler(tra opentracing.Tracer, engine *gin.Engine) {
    apiGroup := engine.Group("/api")

    // 中间件链
    if tra != nil {
        apiGroup.Use(middleware.GinInterceptorWithTrace(tra, false))
    } else {
        apiGroup.Use(middleware.GinInterceptor(true))
    }

    v1API := apiGroup.Group("/v1")
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