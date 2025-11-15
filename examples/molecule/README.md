# go-indigo Molecule Examples

这个文件夹包含了使用 go-indigo molecule 包的各种示例。

## 示例文件

### 1. basic_usage.go

演示基本的分子操作：

- 创建空分子
- 从 SMILES 加载分子
- 克隆分子
- 芳香化/去芳香化
- 氢原子折叠/展开
- 2D 布局
- 归一化和标准化
- 计数连通分量

**运行:**

```bash
cd examples/molecule
go run basic_usage.go
```

### 2. molecule_io.go

演示分子的输入/输出操作：

- 从 SMILES 加载
- 转换为 SMILES（标准和规范）
- 保存到 MOL 文件
- 从 MOL 文件加载
- 转换为 MOL 格式字符串
- JSON 格式的保存和加载
- 从缓冲区加载
- 加载 SMARTS 模式
- 加载查询分子
- 处理混合物

**运行:**

```bash
cd examples/molecule
go run molecule_io.go
```

### 3. molecule_builder.go

演示从头构建分子：

- 构建乙醇（C-C-O）
- 构建苯环
- 构建水分子
- 构建带电荷的分子（乙酸根离子）
- 合并分子

**运行:**

```bash
cd examples/molecule
go run molecule_builder.go
```

### 4. molecule_properties.go

演示分子性质计算：

- 原子和键的计数
- 分子式
- 分子量
- 同位素质量
- 质量组成
- TPSA（拓扑极性表面积）
- 可旋转键数量
- 环的数量
- 自定义属性管理

**运行:**

```bash
cd examples/molecule
go run molecule_properties.go
```

### 5. molecule_inchi.go

演示 InChI 生成和使用：

- 从 SMILES 生成 InChI
- 生成 InChI Key
- 从 InChI 加载分子
- 多种分子的 InChI
- InChI 往返测试
- InChI 警告和日志
- 验证 InChI Key 唯一性

**运行:**

```bash
cd examples/molecule
go run molecule_inchi.go
```

## 环境设置

运行这些示例前，请确保：

### Windows

```cmd
set PATH=%PATH%;D:\for_github\go-indigo\3rd\win
set CGO_ENABLED=1
```

### Linux

```bash
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/path/to/go-indigo/3rd/linux
export CGO_ENABLED=1
```

## 依赖

这些示例需要：

- Go 1.16 或更高版本
- CGO 支持
- Indigo 库（位于 `3rd/win` 或 `3rd/linux`）

## 测试分子

示例中使用的一些测试分子：

| 名称 | SMILES | 描述 |
|------|--------|------|
| 乙醇 | CCO | 简单的醇 |
| 苯 | c1ccccc1 | 芳香环 |
| 乙酸 | CC(=O)O | 羧酸 |
| 阿司匹林 | CC(=O)Oc1ccccc1C(=O)O | 常用药物 |
| 咖啡因 | CN1C=NC2=C1C(=O)N(C(=O)N2C)C | 生物碱 |
| 葡萄糖 | C([C@@H]1[C@H]([C@@H]([C@H](C(O1)O)O)O)O)O | 糖类 |

## 常见问题

### DLL/SO 未找到

确保 `PATH`（Windows）或 `LD_LIBRARY_PATH`（Linux）包含了 Indigo 库的路径。

### CGO 编译错误

确保设置了 `CGO_ENABLED=1` 并且有 C 编译器（如 gcc 或 MSVC）。

### 示例运行失败

检查：

1. Indigo 库文件是否在正确的位置
2. 环境变量是否设置正确
3. Go 模块依赖是否正确安装

## 更多信息

查看 [reaction/SETUP.md](../../reaction/SETUP.md) 了解详细的环境配置说明。
