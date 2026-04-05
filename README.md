# tao-zap

[![Go Report Card](https://goreportcard.com/badge/github.com/taouniverse/tao-zap)](https://goreportcard.com/report/github.com/taouniverse/tao-zap)
[![GoDoc](https://pkg.go.dev/badge/github.com/taouniverse/tao-zap?status.svg)](https://pkg.go.dev/github.com/taouniverse/tao-zap?tab=doc)

Tao Universe Zap 日志组件，集成 uber-go/zap 结构化日志库，支持控制台、文件双输出和日志轮转。

> **设计定位**：tao-zap 是框架级**单例日志基础设施**。它通过 `tao.SetLogger()` + `tao.SetWriter()` 将 zap 日志注册为全局日志后端，所有代码统一通过 `tao.Info()` / `tao.Debug()` 等 API 输出日志，无需也不应直接获取 Logger 实例。

## 安装

```bash
go get github.com/taouniverse/tao-zap
```

## 使用

### 导入

```go
import _ "github.com/taouniverse/tao-zap"
```

### 配置

```yaml
zap:
  logs:
    console:                    # 控制台输出（彩色）
      level: DEBUG
    file:                       # 文件输出（JSON 格式 + lumberjack 轮转）
      level: INFO
      store:
        path: ./logs/app.log
        max_size: 100           # MB
        max_backups: 7
        max_age: 30             # days
        compress: true
        local_time: true
  call_depth: 1                 # 调用栈跳过层数
  coexist: false                # 是否与 tao 内置日志共存（false 则替换默认日志）
  run_after: []
```

### 输出日志

导入后，所有日志自动走 zap 后端，直接使用 tao 框架的日志 API：

```go
tao.Debug("debug message")
tao.Infof("user %s logged in", username)
tao.Warn("something unexpected")
tao.Errorf("failed to connect: %v", err)
```

## 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `logs` | `map[LogType]*LogConfig` | `{console: {level: DEBUG}}` | 日志输出配置，支持 console 和 file 同时开启 |
| `call_depth` | `int` | `1` | 调用栈跳过层数，用于修正日志行号 |
| `coexist` | `bool` | `false` | 是否与 tao 内置日志共存；`false` 时替换默认日志 |

### LogType

| 类型 | 说明 |
|------|------|
| `console` | 控制台输出（彩色文本格式） |
| `file` | 文件输出（JSON 格式 + lumberjack 自动轮转） |

### LogConfig

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `level` | `zapcore.Level` | 日志级别：DEBUG / INFO / WARN / ERROR / PANIC / FATAL |
| `store` | `*Store` | 文件日志专用配置（仅 file 类型需要） |

### Store（文件轮转）

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `path` | `string` | `./test.log` | 日志文件路径 |
| `max_size` | `int` | `1024` | 单文件最大大小（MB） |
| `max_backups` | `int` | `7` | 保留的旧日志文件数 |
| `max_age` | `int` | `30` | 保留天数 |
| `compress` | `bool` | `true` | 是否压缩旧文件 |
| `local_time` | `bool` | `true` | 使用本地时区命名备份文件 |

## 日志级别

| 级别 | 常量 | 说明 |
|------|------|------|
| DEBUG | `DEBUG` | 调试信息 |
| INFO | `INFO` | 一般信息 |
| WARN | `WARNING` | 警告信息 |
| ERROR | `ERROR` | 错误信息 |
| PANIC | `PANIC` | 致命错误（panic） |
| FATAL | `FATAL` | 致命错误（退出） |

## 单元测试

| 测试文件 | 说明 | 运行条件 |
|---------|------|---------|
| `config_test.go` | Config 解析与默认值验证 | 无需外部依赖 |
| `zap_test.go` | Logger 创建与输出测试 | 无需外部服务 |

### 运行单元测试

```bash
# 全部测试可直接运行
go test -v ./...
```

## 依赖

- [uber-go/zap](https://github.com/uber-go/zap) - 高性能结构化日志库
- [lumberjack](https://gopkg.in/natefinch/lumberjack.v2) - 日志文件轮转
