# Go 1.25 升级规划文档

## 项目概况

- **项目**: github.com/yvanz/gin-tmpl
- **原版本**: Go 1.16
- **目标版本**: Go 1.25
- **总文件数**: ~50 个 Go 源文件

## 升级原则 (防止 OOM)

1. **分阶段执行**: 每个阶段独立完成，可随时暂停恢复
2. **小批量修改**: 每次只修改 2-5 个相关文件
3. **先易后难**: 先做安全的替换，再做复杂的重构
4. **独立提交**: 每个阶段完成后单独 commit
5. **类型安全优先**: 使用泛型提高类型安全，而非简单替换

---

## 阶段规划

### 阶段 1: 基础替换 (低风险)
**目标**: `interface{}` → `any`

Go 1.18 引入 `any` 作为 `interface{}` 的别名，这是完全兼容的替换。

**涉及文件**:
1. `internal/common/base.go` - Response.DataSet, CheckParams 等
2. `pkg/gadget/struct.go` - GetTableColumn, StructToMap 等函数参数
3. `pkg/gormdb/repo.go` - CRUD 方法参数
4. `pkg/rediscache/repo.go` - Set 方法参数
5. `pkg/gormdb/basic.go` - 接口定义
6. `pkg/rediscache/basic.go` - 接口定义
7. `pkg/logger/log.go` - 日志函数参数
8. `pkg/middleware/gin-middleware.go` - map[string]interface{}
9. `pkg/rsql/rsql.go` - values 参数

**预计修改**: ~100 处替换
**风险等级**: ★☆☆☆☆ (极低)

---

### 阶段 2: 标准库现代化 (低风险)
**目标**: 使用 Go 1.20+ 标准库改进

**2.1 rand 包自动种子化**
- 文件: `pkg/gadget/uuid.go`
- 改动: 移除 `rand.Seed(time.Now().UnixNano())`
- Go 1.20+ 的 rand 自动种子化

**2.2 slices 包优化**
- 文件: `pkg/gadget/struct.go`
- 已有导入，可以优化:
  - `IsNumber` 函数中的循环查找 → 使用 `slices.Contains`
  - 其他切片操作优化

**2.3 使用 min/max 函数**
- 文件: `pkg/gadget/struct.go`
- Go 1.21+ 引入的内置函数

**预计修改**: ~5-10 处
**风险等级**: ★☆☆☆☆ (极低)

---

### 阶段 3: Context 增强 (中低风险)
**目标**: 使用 `context.WithTimeoutCause` / `WithDeadlineCause`

Go 1.20+ 引入带 cause 的 context 方法，提供更好的错误追踪。

**涉及文件**:
1. `pkg/gadget/tracectx.go` - 可以添加 cause 信息
2. `pkg/kafka/consumer.go` - 如果有 context 使用
3. `pkg/kafka/producer.go` - 如果有 context 使用

**注意**: 需要检查是否会影响现有错误处理逻辑

**预计修改**: ~3-5 处
**风险等级**: ★★☆☆☆ (低)

---

### 阶段 4: CRUD 泛型重构 (中风险)
**目标**: 使用泛型替代 interface{}

这是最复杂的重构，但可以大幅提升类型安全。

**4.1 基础接口泛型化**
- 文件: `pkg/gormdb/basic.go`
- 将接口改为泛型:
```go
// 之前
type GetByIDCrud interface {
    GetByID(model interface{}, id int64) error
}

// 之后
type GetByIDCrud[T any] interface {
    GetByID(model *T, id int64) error
}
```

**4.2 CRUD 实现泛型化**
- 文件: `pkg/gormdb/repo.go`
- 将 CRUDImpl 改为泛型结构体

**4.3 Redis 缓存泛型化**
- 文件: `pkg/rediscache/repo.go`, `pkg/rediscache/basic.go`
- Set 方法可以使用泛型

**4.4 模型层适配**
- 文件: `models/repo/demo.go`, `models/factory/demo.go`
- 适配新的泛型接口

**4.5 业务逻辑层适配**
- 文件: `internal/logic/srvdemo/srv_demo.go`
- 适配新的泛型接口

**注意**: 这是一个破坏性变更，需要级联修改多个文件
**预计修改**: ~50-80 处
**风险等级**: ★★★☆☆ (中等)

---

### 阶段 5: 工具函数泛型化 (中风险)
**目标**: gadget 包泛型优化

**5.1 StructToMap 泛型化**
- 文件: `pkg/gadget/struct.go`
```go
// 之前
func StructToMap(obj interface{}) map[string]interface{}

// 之后
func StructToMap[T any](obj T) map[string]any
```

**5.2 GetTableColumn 泛型化**
- 可以添加泛型约束，限制为 struct 类型

**预计修改**: ~10-15 处
**风险等级**: ★★★☆☆ (中等)

---

### 阶段 6: 响应结构优化 (中低风险)
**目标**: 使用泛型替代 Response.DataSet 的 interface{}

**涉及文件**:
1. `internal/common/base.go` - Response 结构体
2. `internal/handler/demo.go` - 所有 handler

**改动示例**:
```go
// 之前
type Response struct {
    DataSet interface{} `json:"data_set"`
}

// 之后
type Response[T any] struct {
    DataSet T `json:"data_set"`
}
```

**预计修改**: ~20-30 处
**风险等级**: ★★☆☆☆ (低-中)

---

### 阶段 7: 清理与优化 (低风险)
**目标**: 代码清理和现代化

**7.1 移除已弃用的 ioutil 使用**
- 文件: `pkg/middleware/gin-middleware.go`
- Go 1.16 已弃用 ioutil，使用 io 替代

**7.2 使用 errors.Join (Go 1.20+)**
- 检查是否有多个 error 需要合并的场景

**7.3 使用 cmp.Compare (Go 1.21+)**
- 文件: `pkg/gadget/struct.go` 已有 cmp 导入
- 优化比较逻辑

**7.4 其他现代化改进**
- 使用 `clear` 内置函数清空 map (Go 1.21+)
- 使用 `slices` 和 `maps` 标准包函数

**预计修改**: ~10-20 处
**风险等级**: ★☆☆☆☆ (极低)

---

## 执行顺序建议

```
阶段 1 → 阶段 2 → 阶段 3 → 阶段 7 → 阶段 4 → 阶段 5 → 阶段 6
 (基础)  (标准库)  (context) (清理)  (CRUD泛型) (工具泛型) (响应泛型)
```

**理由**:
1. 先做安全的基础替换建立信心
2. 标准库改进是独立的
3. context 增强不影响接口
4. 清理工作可以提前做
5. 泛型重构放在后面，因为它们是结构性变更

---

## 检查清单

每个阶段完成后需要验证:

- [ ] `go build ./...` 编译通过
- [ ] `go vet ./...` 无警告
- [ ] `go test ./...` 测试通过 (如有测试)
- [ ] `go mod tidy` 整理依赖
- [ ] 代码审查关键变更

---

## 回滚策略

每个阶段独立 commit，如果出现问题:

```bash
# 查看最近提交
git log --oneline -10

# 如果需要回滚到阶段 N
git reset --hard <阶段N的commit>

# 或者保留工作区回滚
git reset --soft <阶段N的commit>
```

---

## 预计总工作量

| 阶段 | 预计时间 | 文件数 | 风险 |
|------|---------|--------|------|
| 1    | 30min   | 9      | 极低 |
| 2    | 20min   | 2      | 极低 |
| 3    | 20min   | 3      | 低   |
| 4    | 60min   | 6      | 中等 |
| 5    | 30min   | 1      | 中等 |
| 6    | 40min   | 2      | 低-中|
| 7    | 30min   | 3      | 极低 |
| **总计** | **~4h** | **~26** | - |

---

## 备注

1. 泛型重构是可选的，如果只想做安全升级，可以跳过阶段 4-6
2. 每个阶段建议单独代码审查
3. 如果遇到编译问题，优先检查类型推断和约束
4. Go 1.25 目前处于开发中，确保使用稳定的特性

---

**文档生成时间**: 2026-04-09
**规划版本**: v1.0
