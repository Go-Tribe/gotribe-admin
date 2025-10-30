# 定时任务系统

这是一个优雅的定时任务系统，基于 `github.com/robfig/cron/v3` 实现，提供了完整的任务管理、监控和配置功能。

## 特性

- ✅ **统一接口**: 基于接口的任务定义，易于扩展
- ✅ **配置化**: 支持YAML配置，可动态调整任务参数
- ✅ **监控管理**: 提供任务状态监控和执行历史查询
- ✅ **错误处理**: 支持重试机制和超时控制
- ✅ **生命周期管理**: 优雅的启动和停止机制
- ✅ **API管理**: 提供RESTful API进行任务管理
- ✅ **日志记录**: 完整的任务执行日志和错误追踪

## 架构设计

```
jobs/
├── types.go          # 任务接口和数据结构定义
├── base_job.go       # 基础任务实现
├── manager.go        # 任务管理器
├── registry.go       # 任务注册表
├── config.go         # 任务配置
├── init.go          # 任务初始化
├── sitemap_job.go   # 站点地图生成任务
├── example_job.go   # 示例任务
└── README.md        # 本文档
```

## 快速开始

### 1. 创建自定义任务

```go
package jobs

import (
    "context"
    "time"
    "gotribe-admin/internal/pkg/common"
)

// MyCustomJob 自定义任务
type MyCustomJob struct {
    *BaseJob
}

// NewMyCustomJob 创建自定义任务
func NewMyCustomJob(config JobConfig) *MyCustomJob {
    job := &MyCustomJob{}
    job.BaseJob = NewBaseJob(config, job.execute)
    return job
}

// execute 执行任务逻辑
func (j *MyCustomJob) execute(ctx context.Context) error {
    common.Log.Info("Starting my custom job")

    // 检查上下文是否被取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // 执行具体的业务逻辑
    // ...

    common.Log.Info("My custom job completed")
    return nil
}
```

### 2. 注册任务

```go
// 在 init.go 中添加任务注册
func InitJobs() error {
    // ... 其他任务注册 ...

    // 注册自定义任务
    if config.IsJobEnabled("my_custom") {
        customConfig, _ := config.GetJobConfig("my_custom")
        customJob := NewMyCustomJob(customConfig)
        if err := RegisterJob(customJob); err != nil {
            return err
        }
    }

    return nil
}
```

### 3. 配置任务

在 `config.go` 中添加任务配置：

```go
func DefaultJobsConfig() *JobsConfig {
    return &JobsConfig{
        Enabled: true,
        Jobs: map[string]JobConfig{
            // ... 其他任务配置 ...
            "my_custom": {
                Name:        "my_custom",
                Description: "我的自定义任务",
                Schedule:    "@every 30s",
                Enabled:     true,
                Timeout:     2 * time.Minute,
                RetryCount:  3,
            },
        },
    }
}
```

## API 接口

### 获取任务列表
```http
GET /api/jobs/
```

### 获取任务状态
```http
GET /api/jobs/{name}/status
```

### 获取任务历史
```http
GET /api/jobs/{name}/history?limit=10
```

### 启用任务
```http
POST /api/jobs/{name}/enable
```

### 禁用任务
```http
POST /api/jobs/{name}/disable
```

## 调度表达式

支持标准的cron表达式和预定义表达式：

- `@every 1s` - 每秒执行
- `@every 1m` - 每分钟执行
- `@every 1h` - 每小时执行
- `@daily` - 每天执行
- `@weekly` - 每周执行
- `@monthly` - 每月执行
- `0 0 12 * * *` - 每天12点执行

## 最佳实践

### 1. 任务设计原则
- 任务应该是幂等的
- 支持优雅的中断（context取消）
- 合理设置超时时间
- 实现适当的重试逻辑

### 2. 错误处理
```go
func (j *MyJob) execute(ctx context.Context) error {
    // 使用context检查取消信号
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // 处理业务逻辑
    if err := doSomething(); err != nil {
        // 记录错误但不panic
        common.Log.Errorf("Task failed: %v", err)
        return err
    }

    return nil
}
```

### 3. 资源管理
- 及时释放资源
- 避免长时间占用数据库连接
- 合理使用内存

### 4. 监控和日志
- 记录任务开始和结束时间
- 记录关键业务指标
- 使用结构化日志

## 扩展功能

### 添加任务依赖
可以在任务配置中添加依赖关系，确保任务按顺序执行。

### 分布式任务
可以扩展为支持分布式任务调度，避免单点故障。

### 任务持久化
可以将任务执行历史持久化到数据库，便于长期监控和分析。

## 故障排除

### 常见问题

1. **任务不执行**
   - 检查任务是否启用
   - 验证调度表达式是否正确
   - 查看日志中的错误信息

2. **任务执行失败**
   - 检查任务超时设置
   - 验证重试次数配置
   - 查看具体的错误信息

3. **性能问题**
   - 优化任务执行逻辑
   - 调整任务调度频率
   - 检查资源使用情况

### 调试技巧

1. 启用详细日志
2. 使用任务状态API监控
3. 查看任务执行历史
4. 分析错误模式和频率
