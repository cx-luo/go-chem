# go-chem

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

一个基于 Indigo 库的 Go 化学信息学工具包，使用 CGO 封装提供高性能的分子和反应处理功能。

[English](README_EN.md) | 简体中文

## ✨ 特性

- 🧪 **分子处理**：完整的分子加载、编辑、保存功能
- ⚗️ **反应处理**：化学反应的加载、分析和AAM（原子映射）
- 🎨 **结构渲染**：将分子和反应渲染为图像（PNG、SVG、PDF）
- 🔬 **InChI 支持**：InChI 和 InChIKey 生成与解析
- 📊 **分子属性**：分子量、TPSA、分子式等计算
- 🏗️ **分子构建**：从头构建分子结构
- 🔄 **格式转换**：SMILES、MOL、SDF 等格式互转

## 📦 安装

### 前置要求

1. **Go 1.20+**
2. **Indigo 库**：项目已包含预编译库
   - Windows (x86_64, i386)
   - Linux (x86_64, aarch64)
   - macOS (x86_64, arm64)

### 安装步骤

```bash
# 克隆仓库
git clone https://github.com/cx-luo/go-chem.git
cd go-chem

# 设置环境变量（Windows示例）
set CGO_ENABLED=1
set CGO_CFLAGS=-I%CD%/3rd
set CGO_LDFLAGS=-L%CD%/3rd/windows-x86_64

# 设置环境变量（Linux示例）
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64"
export LD_LIBRARY_PATH=$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH

# 运行测试确认安装成功
go test ./test/molecule/...
```

## 🚀 快速开始

### 加载和渲染分子

```go
package main

import (
   "fmt"
   "github.com/cx-luo/go-chem/core"
   "github.com/cx-luo/go-chem/molecule"
   "github.com/cx-luo/go-chem/render"
)

func main() {
   indigoInit, err := core.IndigoInit()
   if err != nil {
      panic(err)
   }
   
   indigoRender, err := indigoInit.InitRenderer()
   if err != nil {
      fmt.Printf("failed to initialize renderer: %v", err)
   }

   // 从 SMILES 加载分子
   mol, err := indigoInit.LoadMoleculeFromString("c1ccccc1")
   if err != nil {
      panic(err)
   }
   defer mol.Close()

   // 设置渲染选项
   opts := &indigoRender.RenderOptions{
      OutputFormat: "png",
      ImageWidth:   800,
      ImageHeight:  600,
   }
   indigoRender.Options = opts
   indigoRender.Apply()

   // 渲染为 PNG
   indigoRender.RenderToFile(mol.Handle, "benzene.png")
}
```

### 分子属性计算

```go
package main

import (
   "fmt"
   "github.com/cx-luo/go-chem/core"
   "github.com/cx-luo/go-chem/molecule"
)

func main() {
   indigoInit, err := core.IndigoInit()
   if err != nil {
      panic(err)
   }

   // 加载乙醇
   mol, _ := indigoInit.LoadMoleculeFromString("CCO")
   defer mol.Close()

   // 计算分子属性
   mw, _ := mol.MolecularWeight()
   fmt.Printf("分子量: %.2f\n", mw)

   formula, _ := mol.GrossFormula()
   fmt.Printf("分子式: %s\n", formula)

   tpsa, _ := mol.TPSA(false)
   fmt.Printf("TPSA: %.2f\n", tpsa)

   // 转换为 SMILES
   smiles, _ := mol.ToDaylightSmiles()
   fmt.Printf("SMILES: %s\n", smiles)
}
```

### InChI 生成

```go
package main

import (
   "fmt"
   "github.com/cx-luo/go-chem/core"
   "github.com/cx-luo/go-chem/molecule"
)

func main() {
   indigoInit, err := core.IndigoInit()
   if err != nil {
      panic(err)
   }

   indigoInchi, err := indigoInit.InchiInit()
   if err != nil {
      panic(err)
   }

   // 加载分子
   mol, _ := indigoInchi.LoadMoleculeFromString("CC(=O)O")
   defer mol.Close()

   // 生成 InChI
   inchi, _ := indigoInchi.GenerateInChI(mol)
   fmt.Println("InChI:", inchi)

   // 生成 InChIKey
   key, _ := indigoInchi.InchiToKey(inchi)
   fmt.Println("InChIKey:", key)
}
```

### 化学反应处理

```go
package main

import (
   "fmt"
   "github.com/cx-luo/go-chem/core"
   "github.com/cx-luo/go-chem/reaction"
)

func main() {
   indigoInit, err := core.IndigoInit()
   if err != nil {
      panic(err)
   }

   // 加载反应
   rxn, _ := indigoInit.LoadReactionFromString("CCO>>CC=O")
   defer rxn.Close()

   // 获取反应信息
   nReactants, _ := rxn.CountReactants()
   nProducts, _ := rxn.CountProducts()
   fmt.Printf("反应物: %d, 产物: %d\n", nReactants, nProducts)

   // 自动原子映射
   rxn.Automap("discard")

   // 保存为 RXN 文件
   rxn.SaveToFile("reaction.rxn")
}   
```

## 📚 文档

### 核心文档

- [分子处理文档](molecule/README.md) - 分子操作完整指南
- [反应处理文档](reaction/README.md) - 化学反应处理
- [渲染文档](render/README.md) - 结构渲染功能
- [环境设置指南](reaction/SETUP.md) - CGO 环境配置

### 专题文档

- [InChI 实现文档](docs/INCHI.md) - InChI 功能详解
- [API 参考](docs/API.md) - 完整 API 文档
- [示例代码](examples/) - 各种使用示例

## 📂 项目结构

```
go-chem/
├── 3rd/                        # Indigo 预编译库
│   ├── windows-x86_64/         # Windows 64位库
│   ├── windows-i386/           # Windows 32位库
│   ├── linux-x86_64/           # Linux 64位库
│   ├── linux-aarch64/          # Linux ARM64库
│   ├── darwin-x86_64/          # macOS Intel库
│   └── darwin-aarch64/         # macOS Apple Silicon库
├── core/                       # 核心功能
│   ├── indigo.go               # Indigo 库核心功能
│   ├── indigo_helper.go        # Indigo 辅助功能
│   ├── indigo_inchi.go         # Indigo InChI 功能
│   ├── indigo_molecule.go      # Indigo 分子功能
│   └── indigo_reaction.go      # Indigo 反应功能
├── molecule/                   # 分子处理包
│   ├── README.md               # 分子处理文档
│   ├── molecule.go             # 核心分子结构
│   ├── molecule_atom.go        # 原子操作
│   ├── molecule_builder.go     # 分子构建
│   ├── molecule_match.go       # 分子匹配
│   ├── molecule_properties.go  # 属性计算
│   └── molecule_saver.go       # 分子保存
├── reaction/                   # 反应处理包
│   ├── README.md               # 反应处理文档
│   ├── reaction.go             # 核心反应结构
│   ├── reaction_automap.go     # 自动原子映射
│   ├── reaction_helpers.go     # 反应辅助功能
│   ├── reaction_iterator.go    # 反应迭代器
│   ├── reaction_loader.go      # 反应加载
│   └── reaction_saver.go       # 反应保存
├── render/                     # 渲染包
│   ├── README.md               # 渲染文档
│   └── render.go               # 渲染功能
├── test/                       # 测试文件
│   ├── molecule/               # 分子测试
│   ├── reaction/               # 反应测试
│   └── render/                 # 渲染测试
├── examples/                   # 示例代码
│   ├── molecule/               # 分子示例
│   ├── reaction/               # 反应示例
│   └── render/                 # 渲染示例
├── docs/                       # 文档
└── README.md                   # 本文件
```

## 🔧 支持的功能

### 分子操作

- ✅ 从 SMILES、MOL、SDF 加载分子
- ✅ 保存为 MOL、SMILES、JSON 格式
- ✅ 分子属性计算（分子量、TPSA、分子式等）
- ✅ 原子和键的添加、删除、修改
- ✅ 芳香化和去芳香化
- ✅ 氢原子折叠和展开
- ✅ 2D 布局和清理
- ✅ 分子标准化和归一化

### 反应操作

- ✅ 从 Reaction SMILES、RXN 文件加载
- ✅ 保存为 RXN 文件
- ✅ 添加反应物、产物、催化剂
- ✅ 自动原子到原子映射（AAM）
- ✅ 反应中心检测
- ✅ 反应组件迭代

### 渲染功能

- ✅ PNG、SVG、PDF 输出
- ✅ 自定义图像大小和样式
- ✅ 网格渲染（多个分子）
- ✅ 参考原子对齐
- ✅ 立体化学显示
- ✅ 原子/键标签显示

### InChI 支持

- ✅ 标准 InChI 生成
- ✅ InChIKey 生成
- ✅ 从 InChI 加载分子
- ✅ 警告和日志信息
- ✅ 辅助信息输出

## 🧪 测试

```bash
# 运行所有测试
go test ./test/...

# 运行特定包的测试
go test ./test/molecule/...
go test ./test/reaction/...
go test ./test/render/...

# 运行带详细输出的测试
go test -v ./test/...

# 运行特定测试
go test ./test/molecule/ -run TestLoadMoleculeFromString
```

## 📊 性能

- 基于 C++ Indigo 库，性能优秀
- CGO 调用开销最小化
- 内存自动管理（使用 runtime.SetFinalizer）
- 支持大规模分子处理

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request 或创建 Issue。

### 开发环境设置

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 Apache License 2.0 许可证。详见 [LICENSE](LICENSE) 文件。

### 第三方许可

- **Indigo Toolkit**: Apache License 2.0
- Copyright © 2009-Present EPAM Systems

## 🙏 致谢

- [EPAM Indigo](https://github.com/epam/Indigo) - 优秀的化学信息学工具包
- 所有贡献者和使用者

## 📮 联系方式

- 作者：chengxiang.luo
- 邮箱：<chengxiang.luo@foxmail.com>
- GitHub：[@cx-luo](https://github.com/cx-luo)

## 🔗 相关链接

- [Indigo 官方文档](https://lifescience.opensource.epam.com/indigo/)
- [Go 官方文档](https://golang.org/doc/)
- [CGO 文档](https://golang.org/cmd/cgo/)

---

⭐ 如果这个项目对你有帮助，请给一个 Star！
