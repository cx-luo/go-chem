# InChI CGO 集成指南

本文档说明如何使用 CGO 调用 InChI 动态链接库，以实现高性能和标准兼容的 InChI 生成。

## 概述

本项目提供两种 InChI 生成方式：

1. **Pure Go 实现** (`molecule_inchi.go`): 纯 Go 语言实现，无需外部依赖
2. **CGO 绑定** (`molecule_inchi_cgo.go`): 调用官方 InChI 库，保证标准兼容性和性能

## CGO 版本的优势

### 为什么选择 CGO 版本？

1. **标准兼容性**: 使用 IUPAC 官方 InChI 库，保证 100% 兼容性
2. **性能优势**: C 库经过高度优化，性能更好
3. **完整功能**: 支持所有 InChI 选项和特性
4. **可靠性**: 经过广泛测试和验证的库

### 对比

| 特性 | Pure Go | CGO |
|------|---------|-----|
| 外部依赖 | 无 | 需要 InChI 库 |
| 性能 | 中等 | 高 |
| 标准兼容性 | 部分 | 100% |
| 跨平台编译 | 简单 | 需要库文件 |
| 功能完整性 | 基本功能 | 完整功能 |
| 部署 | 简单 | 需要动态库 |

## 项目结构

```
go-chem/
├── 3rd/                        # 第三方库目录
│   ├── inchi_api.h            # InChI API 头文件
│   ├── libinchi.dll           # Windows 动态库
│   └── libinchi.so            # Linux 动态库
├── molecule/
│   ├── molecule_inchi.go      # Pure Go 实现
│   └── molecule_inchi_cgo.go  # CGO 绑定实现
└── examples/
    ├── inchi_example.go       # Pure Go 示例
    └── inchi_cgo_example.go   # CGO 示例
```

## 安装和配置

### 1. 获取 InChI 库

InChI 动态库文件已包含在 `3rd/` 目录中：

- **Windows**: `3rd/libinchi.dll`
- **Linux**: `3rd/libinchi.so`

如果需要最新版本，可从 [InChI Trust](https://www.inchi-trust.org/downloads/) 下载。

### 2. 环境配置

#### Windows

```bash
# 确保 libinchi.dll 在系统 PATH 或项目 3rd/ 目录
set PATH=%PATH%;%CD%\3rd
```

#### Linux

```bash
# 设置库路径
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(pwd)/3rd

# 或者将库复制到系统目录
sudo cp 3rd/libinchi.so /usr/local/lib/
sudo ldconfig
```

### 3. 构建项目

#### 使用 CGO 构建

```bash
# 标准构建（启用 CGO）
go build -o inchi_cgo examples/inchi_cgo_example.go

# 显式启用 CGO
CGO_ENABLED=1 go build examples/inchi_cgo_example.go

# Windows 静态链接（可选）
go build -ldflags="-extldflags=-static" examples/inchi_cgo_example.go
```

#### 不使用 CGO 构建（Pure Go）

```bash
# 禁用 CGO
CGO_ENABLED=0 go build examples/inchi_example.go
```

### 4. 运行示例

```bash
# Windows
.\inchi_cgo.exe

# Linux
./inchi_cgo
```

## 使用方法

### 基本用法

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/cx-luo/go-chem/molecule"
)

func main() {
    // 创建 CGO 生成器
    generator := molecule.NewInChIGeneratorCGO()
    
    // 解析分子
    loader := molecule.SmilesLoader{}
    mol, err := loader.Parse("CCO") // 乙醇
    if err != nil {
        log.Fatal(err)
    }
    
    // 生成 InChI
    result, err := generator.GenerateInChI(mol)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("InChI:", result.InChI)
    fmt.Println("InChIKey:", result.InChIKey)
}
```

### 设置选项

```go
generator := molecule.NewInChIGeneratorCGO()

// 标准 InChI（默认）
generator.SetOptions("")

// 固定氢层
generator.SetOptions("FixedH")

// 重连接金属
generator.SetOptions("RecMet")

// 多个选项
generator.SetOptions("FixedH RecMet")
```

### 常用选项

| 选项 | 说明 |
|------|------|
| `""` | 标准 InChI（默认） |
| `FixedH` | 包含固定氢层 |
| `RecMet` | 重连接金属 |
| `SNon` | 包含省略的未定义/未知立体化学 |
| `AuxNone` | 不生成辅助信息 |
| `Wnumber` | 包含警告编号 |

完整选项列表见 [InChI 技术手册](https://www.inchi-trust.org/downloads/)。

### 仅生成 InChIKey

```go
import "github.com/cx-luo/go-chem/molecule"

// 从 InChI 生成 InChIKey
inchi := "InChI=1S/CH4/h1H4"
key, err := molecule.GenerateInChIKeyCGO(inchi)
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChIKey:", key)
```

### 获取库版本

```go
version := molecule.GetInChIVersion()
fmt.Println("InChI Library Version:", version)
```

## CGO 实现细节

### 架构

CGO 实现遵循 Indigo 的 `inchi_wrapper.cpp` 架构：

```
Go 应用
  ↓
molecule_inchi_cgo.go (Go + CGO)
  ↓
inchi_api.h (C 头文件)
  ↓
libinchi.dll/so (InChI 库)
```

### 主要函数

#### 1. 创建 InChI 输入

```go
func (g *InChIGeneratorCGO) createInChIInput(mol *Molecule) (*C.inchi_Input, error)
```

转换 Go 分子结构到 C 的 `inchi_Input` 结构：
- 原子数组：元素、坐标、电荷、自由基、同位素
- 键数组：键类型、立体化学
- 立体元素：顺反异构、四面体

参考: `inchi_wrapper.cpp`, `generateInchiInput` (406-611 行)

#### 2. 生成 InChI

```go
func (g *InChIGeneratorCGO) GenerateInChI(mol *Molecule) (*InChIResult, error)
```

调用 InChI 库的 `GetINCHI()` 函数生成 InChI。

参考: `inchi_wrapper.cpp`, `saveMoleculeIntoInchi` (620-703 行)

#### 3. 生成 InChIKey

```go
func GenerateInChIKeyCGO(inchi string) (string, error)
```

调用 InChI 库的 `GetINCHIKeyFromINCHI()` 函数。

参考: `inchi_wrapper.cpp`, `InChIKey` (705-730 行)

### C 辅助函数

为了简化 CGO 调用，实现了以下 C 辅助函数：

```c
// 分配/释放结构
inchi_Input* alloc_inchi_input();
void free_inchi_input(inchi_Input* inp);
inchi_Atom* alloc_atoms(int count);
inchi_Stereo0D* alloc_stereo(int count);

// 设置数据
void set_atom_data(...);
void add_bond(...);
void set_hydrogen_count(...);
void set_stereo_data(...);
```

### 内存管理

CGO 代码严格管理内存：

1. **Go 管理的内存**: 使用 `runtime.SetFinalizer`
2. **C 管理的内存**: 使用 `defer` 确保释放
3. **InChI 库分配的内存**: 调用 `FreeINCHI()` 和 `FreeStructFromINCHI()`

```go
// 示例
inp := C.alloc_inchi_input()
defer C.free_inchi_input(inp)

var out C.inchi_Output
defer C.FreeINCHI(&out)
```

## 测试

### 运行 CGO 测试

```bash
# 运行所有测试
CGO_ENABLED=1 go test -v ./test -run TestInChICGO

# 运行特定测试
go test -v ./test -run TestInChICGO_Basic
```

### 基准测试

```bash
# 对比 CGO vs Pure Go 性能
go test -bench=BenchmarkInChI -benchmem ./test
```

### 验证结果

```bash
# 使用官方 InChI 工具验证
./inchi -STDIO < molecule.mol
```

## 部署

### 独立可执行文件（Windows）

```bash
# 静态链接 InChI 库
go build -ldflags="-extldflags=-static" -o app.exe

# 或者将 libinchi.dll 与可执行文件打包
cp 3rd/libinchi.dll ./
```

### 容器部署（Docker）

```dockerfile
FROM golang:1.21 AS builder

# 安装 CGO 依赖
RUN apt-get update && apt-get install -y gcc

# 复制源码
WORKDIR /app
COPY .. .

# 构建
RUN CGO_ENABLED=1 go build -o myapp

# 运行时镜像
FROM debian:bookworm-slim
COPY --from=builder /app/myapp /
COPY --from=builder /app/3rd/libinchi.so /usr/local/lib/
RUN ldconfig

ENTRYPOINT ["/myapp"]
```

### Linux 发行版

```bash
# 方式 1: 系统库
sudo cp 3rd/libinchi.so /usr/local/lib/
sudo ldconfig

# 方式 2: rpath
go build -ldflags="-r /path/to/libs"

# 方式 3: 环境变量
export LD_LIBRARY_PATH=/path/to/3rd:$LD_LIBRARY_PATH
```

## 故障排除

### 常见问题

#### 1. "cannot find -linchi"

**原因**: 找不到 InChI 库文件

**解决**:
```bash
# 检查库文件
ls -la 3rd/

# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd:$LD_LIBRARY_PATH

# Windows
set PATH=%PATH%;%CD%\3rd
```

#### 2. "undefined reference to GetINCHI"

**原因**: CGO 链接配置错误

**解决**:
```bash
# 检查 CGO 配置
go env CGO_ENABLED
go env CGO_LDFLAGS

# 重新构建
go clean -cache
CGO_ENABLED=1 go build
```

#### 3. "InChI generation failed"

**原因**: 分子结构不兼容

**解决**:
- 检查分子是否包含伪原子或 R-基团
- 查看错误消息和警告
- 尝试简化分子结构

#### 4. 运行时找不到库

**Linux**:
```bash
# 临时
export LD_LIBRARY_PATH=/path/to/3rd:$LD_LIBRARY_PATH

# 永久（添加到 ~/.bashrc）
echo 'export LD_LIBRARY_PATH=/path/to/3rd:$LD_LIBRARY_PATH' >> ~/.bashrc
```

**Windows**:
```bash
# 将 DLL 放在可执行文件同目录
# 或添加到系统 PATH
```

### 调试技巧

#### 启用详细日志

```go
result, err := generator.GenerateInChI(mol)
if err != nil {
    log.Printf("Error: %v\n", err)
}

// 检查警告和日志
for _, warning := range result.Warnings {
    log.Printf("Warning: %s\n", warning)
}
for _, logMsg := range result.Log {
    log.Printf("Log: %s\n", logMsg)
}
```

#### 验证库加载

```go
version := molecule.GetInChIVersion()
if version == "unknown" {
    log.Fatal("Failed to load InChI library")
}
log.Printf("InChI library version: %s\n", version)
```

## 性能优化

### 批量处理

```go
generator := molecule.NewInChIGeneratorCGO()
loader := molecule.SmilesLoader{}

// 复用生成器和加载器
for _, smiles := range smilesList {
    mol, _ := loader.Parse(smiles)
    result, _ := generator.GenerateInChI(mol)
    // 处理结果...
}
```

### 并发处理

```go
import "sync"

var wg sync.WaitGroup
results := make(chan *molecule.InChIResult, len(smilesList))

for _, smiles := range smilesList {
    wg.Add(1)
    go func(s string) {
        defer wg.Done()
        
        // 每个 goroutine 使用独立的生成器
        gen := molecule.NewInChIGeneratorCGO()
        loader := molecule.SmilesLoader{}
        
        mol, _ := loader.Parse(s)
        result, _ := gen.GenerateInChI(mol)
        results <- result
    }(smiles)
}

wg.Wait()
close(results)
```

**注意**: InChI 库内部可能使用互斥锁，并发收益有限。

## 参考资料

### InChI 文档

- [InChI Trust 官网](https://www.inchi-trust.org/)
- [InChI 技术手册](https://www.inchi-trust.org/downloads/)
- [InChI API 文档](https://www.inchi-trust.org/downloads/inchi-api/)

### 实现参考

- **Indigo C++ 源码**: `indigo-core/molecule/src/inchi_wrapper.cpp`
- **本项目文档**: `INCHI_IMPLEMENTATION.md`

### CGO 文档

- [Go CGO 官方文档](https://golang.org/cmd/cgo/)
- [CGO 最佳实践](https://github.com/golang/go/wiki/cgo)

## 总结

### 何时使用 CGO 版本？

✅ **推荐使用 CGO**:
- 需要标准兼容的 InChI
- 处理复杂分子和立体化学
- 性能敏感的应用
- 需要所有 InChI 选项

❌ **不推荐使用 CGO**:
- 简单部署需求
- 跨平台编译困难
- 不需要完整功能

### 何时使用 Pure Go 版本？

✅ **推荐使用 Pure Go**:
- 简单部署
- 跨平台编译
- 基本 InChI 功能足够
- 不想依赖外部库

选择适合您需求的版本！

