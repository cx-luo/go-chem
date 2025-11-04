# go-chem 示例代码

本目录包含 go-chem 各个模块的完整使用示例。

## 📂 示例目录

### Molecule 示例

详细的分子处理示例，请查看 [molecule/README.md](molecule/README.md)

- **[basic_usage.go](molecule/basic_usage.go)** - 分子基本操作
- **[molecule_io.go](molecule/molecule_io.go)** - 分子输入输出
- **[molecule_builder.go](molecule/molecule_builder.go)** - 分子构建
- **[molecule_properties.go](molecule/molecule_properties.go)** - 分子属性计算
- **[molecule_inchi.go](molecule/molecule_inchi.go)** - InChI 生成和使用

### Reaction 示例

- **[example_reaction.go](example_reaction.go)** - 化学反应处理完整示例

包含内容：

- 加载反应（SMILES, RXN 文件）
- 反应组件迭代
- 自动原子映射（AAM）
- 反应保存
- 反应标准化

**运行:**

```bash
go run example_reaction.go
```

### Render 示例

- **[example_render.go](example_render.go)** - 结构渲染完整示例

包含内容：

- 基本分子渲染
- 反应渲染
- 网格渲染
- 自定义渲染选项
- 多种输出格式（PNG、SVG、PDF）
- 内存缓冲区渲染

**运行:**

```bash
go run example_render.go
```

### InChI 示例

- **[inchi_usage.go](inchi_usage.go)** - InChI 功能使用示例

包含内容：

- InChI 生成
- InChIKey 生成
- 从 InChI 加载分子
- 批量处理
- InChI 往返转换

**运行:**

```bash
go run inchi_usage.go
```

## 🚀 快速开始

### 环境准备

#### Windows

```cmd
REM 设置环境变量
set CGO_ENABLED=1
set CGO_CFLAGS=-ID:\path\to\go-chem\3rd
set CGO_LDFLAGS=-LD:\path\to\go-chem\3rd\windows-x86_64
set PATH=%PATH%;D:\path\to\go-chem\3rd\windows-x86_64

REM 运行示例
cd examples
go run example_reaction.go
```

#### Linux

```bash
# 设置环境变量
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64"
export LD_LIBRARY_PATH="$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH"

# 运行示例
cd examples
go run example_reaction.go
```

#### macOS

```bash
# 设置环境变量
export CGO_ENABLED=1
export CGO_CFLAGS="-I$(pwd)/3rd"

# 根据架构选择
export CGO_LDFLAGS="-L$(pwd)/3rd/darwin-aarch64"  # M1/M2
# export CGO_LDFLAGS="-L$(pwd)/3rd/darwin-x86_64"  # Intel

export DYLD_LIBRARY_PATH="$(pwd)/3rd/darwin-aarch64:$DYLD_LIBRARY_PATH"

# 运行示例
cd examples
go run example_reaction.go
```

### 运行所有示例

```bash
# 分子示例
cd molecule
go run basic_usage.go
go run molecule_io.go
go run molecule_builder.go
go run molecule_properties.go
go run molecule_inchi.go

# 反应示例
cd ..
go run example_reaction.go

# 渲染示例
go run example_render.go

# InChI 示例
go run inchi_usage.go
```

## 📚 示例分类

### 入门级示例

适合初学者：

1. **[molecule/basic_usage.go](molecule/basic_usage.go)** - 从这里开始
2. **[example_reaction.go](example_reaction.go)** - 学习反应处理
3. **[example_render.go](example_render.go)** - 学习结构渲染

### 进阶示例

适合有一定基础的用户：

1. **[molecule/molecule_builder.go](molecule/molecule_builder.go)** - 从头构建分子
2. **[molecule/molecule_properties.go](molecule/molecule_properties.go)** - 复杂属性计算
3. **[molecule/molecule_inchi.go](molecule/molecule_inchi.go)** - InChI 高级用法

### 专题示例

特定功能的深入示例：

1. **[inchi_usage.go](inchi_usage.go)** - InChI 完整功能
2. **[molecule/molecule_io.go](molecule/molecule_io.go)** - 文件格式转换

## 🎯 按功能索引

### 分子加载

- SMILES 加载：`molecule/basic_usage.go`
- MOL 文件加载：`molecule/molecule_io.go`
- InChI 加载：`molecule/molecule_inchi.go`
- 查询分子加载：`molecule/molecule_io.go`

### 分子保存

- SMILES 输出：`molecule/molecule_io.go`
- MOL 文件保存：`molecule/molecule_io.go`
- JSON 保存：`molecule/molecule_io.go`

### 分子构建

- 添加原子和键：`molecule/molecule_builder.go`
- 构建环结构：`molecule/molecule_builder.go`
- 合并分子：`molecule/molecule_builder.go`

### 分子属性

- 分子量：`molecule/molecule_properties.go`
- 分子式：`molecule/molecule_properties.go`
- TPSA：`molecule/molecule_properties.go`
- 自定义属性：`molecule/molecule_properties.go`

### 反应处理

- 反应加载：`example_reaction.go`
- 原子映射：`example_reaction.go`
- 反应迭代：`example_reaction.go`
- 反应保存：`example_reaction.go`

### 渲染功能

- 分子渲染：`example_render.go`
- 反应渲染：`example_render.go`
- 网格渲染：`example_render.go`
- 自定义样式：`example_render.go`

### InChI 功能

- InChI 生成：`inchi_usage.go`, `molecule/molecule_inchi.go`
- InChIKey 生成：`inchi_usage.go`
- InChI 转分子：`molecule/molecule_inchi.go`

## 🔧 常见问题

### 运行时错误

#### Windows: `exit status 0xc0000135`

DLL 文件未找到。解决方案：

```cmd
set PATH=%PATH%;D:\path\to\go-chem\3rd\windows-x86_64
```

#### Linux: `error while loading shared libraries`

共享库未找到。解决方案：

```bash
export LD_LIBRARY_PATH="$(pwd)/3rd/linux-x86_64:$LD_LIBRARY_PATH"
```

#### macOS: `Library not loaded`

动态库未找到。解决方案：

```bash
export DYLD_LIBRARY_PATH="$(pwd)/3rd/darwin-aarch64:$DYLD_LIBRARY_PATH"
```

### 编译错误

#### `indigo.h: No such file or directory`

解决方案：

```bash
export CGO_CFLAGS="-I$(pwd)/3rd"
```

#### `undefined reference to 'indigoCreateMolecule'`

解决方案：

```bash
export CGO_LDFLAGS="-L$(pwd)/3rd/linux-x86_64 -lindigo"
```

## 📖 学习路径

### 第1天：基础入门

1. 阅读 [molecule/basic_usage.go](molecule/basic_usage.go)
2. 运行并理解每个函数
3. 尝试修改 SMILES 输入

### 第2天：输入输出

1. 学习 [molecule/molecule_io.go](molecule/molecule_io.go)
2. 尝试不同的文件格式
3. 创建自己的 MOL 文件

### 第3天：分子构建

1. 研究 [molecule/molecule_builder.go](molecule/molecule_builder.go)
2. 从头构建一个简单分子
3. 尝试构建更复杂的结构

### 第4天：属性计算

1. 学习 [molecule/molecule_properties.go](molecule/molecule_properties.go)
2. 计算不同分子的属性
3. 理解 TPSA 和其他药物性质

### 第5天：反应和渲染

1. 运行 [example_reaction.go](example_reaction.go)
2. 运行 [example_render.go](example_render.go)
3. 创建自己的反应示例

### 第6天：InChI 高级用法

1. 深入 [molecule/molecule_inchi.go](molecule/molecule_inchi.go)
2. 理解 InChI 的层次结构
3. 使用 InChIKey 进行分子比较

### 第7天：综合应用

1. 结合多个示例创建完整应用
2. 批量处理分子
3. 生成分子数据库

## 💡 提示和技巧

### 性能优化

```go
// ✅ 好的做法：重用对象
mol, _ := molecule.LoadMoleculeFromString("CCO")
defer mol.Close()

for i := 0; i < 100; i++ {
    clone, _ := mol.Clone()
    // 处理克隆
    clone.Close()
}

// ❌ 避免：重复加载
for i := 0; i < 100; i++ {
    mol, _ := molecule.LoadMoleculeFromString("CCO")  // 慢
    mol.Close()
}
```

### 错误处理

```go
// ✅ 总是检查错误
mol, err := molecule.LoadMoleculeFromString(smiles)
if err != nil {
    log.Fatalf("加载失败: %v", err)
}
defer mol.Close()

// ❌ 不要忽略错误
mol, _ := molecule.LoadMoleculeFromString(smiles)
```

### 资源管理

```go
// ✅ 使用 defer 确保资源释放
func ProcessMolecule(smiles string) error {
    mol, err := molecule.LoadMoleculeFromString(smiles)
    if err != nil {
        return err
    }
    defer mol.Close()  // 总是会执行

    // 处理分子...
    return nil
}
```

## 📞 获取帮助

- 查看 [API 文档](../docs/API.md)
- 阅读 [常见问题](../docs/FAQ.md)
- 在 GitHub 创建 Issue
- 发送邮件至 <chengxiang.luo@foxmail.com>

## 📄 许可证

所有示例代码采用 Apache License 2.0 许可证。

---

⭐ **提示**: 建议按顺序学习示例，从简单到复杂！
