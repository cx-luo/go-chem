# Go-Chem InChI 实现

本项目提供两种 InChI (International Chemical Identifier) 和 InChIKey 生成实现。

## 🎯 快速开始

### 方式 1: Pure Go 实现（推荐用于简单场景）

**优点**: 无需外部依赖，跨平台编译简单

```go
import "github.com/cx-luo/go-chem/molecule"

// 从 SMILES 生成 InChI
result, err := molecule.GetInChIFromSMILES("CCO") // 乙醇
if err != nil {
    log.Fatal(err)
}

fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

**运行示例**:
```bash
go run examples/inchi_example.go
```

### 方式 2: CGO 实现（推荐用于生产环境）

**优点**: 使用官方 InChI 库，100% 标准兼容，性能更好

```go
import "github.com/cx-luo/go-chem/molecule"

// 创建 CGO 生成器
generator := molecule.NewInChIGeneratorCGO()

// 解析分子
loader := molecule.SmilesLoader{}
mol, _ := loader.Parse("CCO")

// 生成 InChI
result, _ := generator.GenerateInChI(mol)
fmt.Println("InChI:", result.InChI)
fmt.Println("InChIKey:", result.InChIKey)
```

**运行示例**:
```bash
# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd:$LD_LIBRARY_PATH
CGO_ENABLED=1 go run examples/inchi_cgo_example.go

# Windows
set PATH=%PATH%;%CD%\3rd
go run examples/inchi_cgo_example.go
```

## 📊 功能对比

| 特性 | Pure Go | CGO |
|------|---------|-----|
| **外部依赖** | ✅ 无 | ⚠️ 需要 libinchi |
| **标准兼容性** | ⚠️ 部分 | ✅ 100% |
| **性能** | ⚠️ 中等 | ✅ 优秀 |
| **跨平台编译** | ✅ 简单 | ⚠️ 需要库文件 |
| **功能完整性** | ⚠️ 基本 | ✅ 完整 |
| **部署难度** | ✅ 简单 | ⚠️ 中等 |

## 📚 详细文档

- **[INCHI_IMPLEMENTATION.md](INCHI_IMPLEMENTATION.md)** - Pure Go 实现详解
  - InChI 算法原理
  - 与 Indigo C++ 代码对应关系
  - 各层实现细节

- **[INCHI_SUMMARY.md](INCHI_SUMMARY.md)** - 中文摘要
  - 主要功能概述
  - 使用方法
  - 改进点列表

- **[INCHI_CGO_GUIDE.md](INCHI_CGO_GUIDE.md)** - CGO 集成指南
  - CGO 环境配置
  - 构建和部署
  - 性能优化
  - 故障排除

## 🚀 安装

```bash
# 下载项目
git clone https://github.com/cx-luo/go-chem.git
cd go-chem

# Pure Go 版本（无需额外步骤）
go run examples/inchi_example.go

# CGO 版本（需要配置库路径）
# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd:$LD_LIBRARY_PATH
CGO_ENABLED=1 go run examples/inchi_cgo_example.go

# Windows
set PATH=%PATH%;%CD%\3rd
go run examples/inchi_cgo_example.go
```

## 💡 使用示例

### 1. 基本用法

```go
// Pure Go
result, err := molecule.GetInChIFromSMILES("c1ccccc1") // 苯
// InChI: InChI=1S/C6H6/c1-2-4-6-5-3-1/h1-6H

// CGO
generator := molecule.NewInChIGeneratorCGO()
loader := molecule.SmilesLoader{}
mol, _ := loader.Parse("c1ccccc1")
result, _ := generator.GenerateInChI(mol)
```

### 2. InChIKey 生成

```go
// Pure Go
key, _ := molecule.GenerateInChIKey("InChI=1S/CH4/h1H4")

// CGO
key, _ := molecule.GenerateInChIKeyCGO("InChI=1S/CH4/h1H4")
```

### 3. 自定义选项（仅 Pure Go）

```go
generator := molecule.NewInChIGenerator()
generator.SetOptions(molecule.InChIOptions{
    FixedH:  true,  // 包含氢层
    RecMet:  false,
    AuxInfo: false,
})
```

### 4. 自定义选项（CGO）

```go
generator := molecule.NewInChIGeneratorCGO()
generator.SetOptions("FixedH RecMet")
```

## 🔬 测试

```bash
# 运行所有 InChI 测试
go test -v ./test -run TestInChI

# 运行 CGO 测试
CGO_ENABLED=1 go test -v ./test -run TestInChICGO

# 性能基准测试
go test -bench=BenchmarkInChI -benchmem ./test
```

## 🏗️ 项目结构

```
go-chem/
├── 3rd/                          # 第三方库
│   ├── inchi_api.h              # InChI API 头文件
│   ├── libinchi.dll             # Windows 动态库
│   └── libinchi.so              # Linux 动态库
├── molecule/
│   ├── molecule_inchi.go        # Pure Go 实现
│   └── molecule_inchi_cgo.go    # CGO 实现
├── examples/
│   ├── inchi_example.go         # Pure Go 示例
│   └── inchi_cgo_example.go     # CGO 示例
├── test/
│   ├── inchi_test.go            # Pure Go 测试
│   └── inchi_stereochemistry_test.go
├── indigo-core/                 # C++ 参考实现
│   └── molecule/src/
│       ├── inchi_wrapper.cpp
│       ├── molecule_inchi.cpp
│       └── molecule_inchi_layers.cpp
└── docs/
    ├── INCHI_IMPLEMENTATION.md  # Pure Go 实现文档
    ├── INCHI_SUMMARY.md         # 中文摘要
    └── INCHI_CGO_GUIDE.md       # CGO 指南
```

## 📖 InChI 层次结构

InChI 由多个层组成，每层提供特定信息：

```
InChI=1S/C6H12O6/c7-1-2-3(8)4(9)5(10)6(11)12-2/h2-11H,1H2/t2-,3-,4+,5-,6?/m1/s1
│       │       │                              │          │   │
│       │       │                              │          │   └─ 立体类型 (/s)
│       │       │                              │          └─ 对映体 (/m)
│       │       │                              └─ 四面体立体化学 (/t)
│       │       └─ 氢原子 (/h)
│       └─ 连接表 (/c)
└─ 化学式
```

## 🎓 算法参考

本实现基于以下资源：

1. **Indigo 开源项目**
   - `indigo-core/molecule/src/inchi_wrapper.cpp`
   - `indigo-core/molecule/src/molecule_inchi.cpp`
   - `indigo-core/molecule/src/molecule_inchi_layers.cpp`

2. **IUPAC InChI 规范**
   - [InChI Trust 官网](https://www.inchi-trust.org/)
   - [技术手册](https://www.inchi-trust.org/downloads/)

3. **学术文献**
   - Goodman, J.M., et al. (2012) "InChI version 1, three years on"
   - Heller, S., et al. (2013) "InChI - the worldwide chemical structure identifier standard"

## 🛠️ 构建说明

### Pure Go 版本

```bash
# 标准构建
go build -o inchi_pure examples/inchi_example.go

# 禁用 CGO（确保纯 Go）
CGO_ENABLED=0 go build examples/inchi_example.go

# 跨平台编译
GOOS=linux GOARCH=amd64 go build examples/inchi_example.go
GOOS=windows GOARCH=amd64 go build examples/inchi_example.go
```

### CGO 版本

```bash
# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd:$LD_LIBRARY_PATH
CGO_ENABLED=1 go build -o inchi_cgo examples/inchi_cgo_example.go

# Windows
set PATH=%PATH%;%CD%\3rd
go build -o inchi_cgo.exe examples/inchi_cgo_example.go

# 静态链接（Windows）
go build -ldflags="-extldflags=-static" examples/inchi_cgo_example.go
```

## 🐛 故障排除

### CGO 编译错误

```bash
# 检查 CGO 是否启用
go env CGO_ENABLED

# 启用 CGO
export CGO_ENABLED=1  # Linux/Mac
set CGO_ENABLED=1     # Windows

# 清理缓存
go clean -cache
go build
```

### 运行时找不到库

```bash
# Linux
export LD_LIBRARY_PATH=$(pwd)/3rd:$LD_LIBRARY_PATH
./inchi_cgo

# Windows
# 将 libinchi.dll 放在可执行文件同目录
copy 3rd\libinchi.dll .
inchi_cgo.exe
```

## 🤝 贡献

欢迎贡献代码改进！特别是：

- ✅ 规范化编号算法改进
- ✅ 立体化学处理增强
- ✅ 性能优化
- ✅ InChI 解析功能实现
- ✅ 更多测试用例

## 📝 许可证

Apache License 2.0

## 📧 联系

如有问题或建议，请提交 Issue 或 Pull Request。

---

**选择建议**:

- 🚀 **快速原型**: 使用 Pure Go 版本
- 🏭 **生产环境**: 使用 CGO 版本
- 📦 **简单部署**: 使用 Pure Go 版本
- 🎯 **标准兼容**: 使用 CGO 版本

