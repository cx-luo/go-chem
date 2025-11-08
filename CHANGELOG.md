# 更新日志 / Changelog

本文档记录 go-chem 项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 计划中的功能

- 分子指纹和相似度计算
- 子结构搜索
- SMARTS 模式匹配
- 反应查询和搜索
- 3D 坐标生成
- 力场能量计算

## [0.4.4] - 2025-11-08

### 新增

- **Reaction 包增强**:
  - `GetReactantMolecule()` - 获取单个反应物分子
  - `GetProductMolecule()` - 获取单个产物分子
  - `GetCatalystMolecule()` - 获取单个催化剂分子
  - `GetReactantMolecules()` - 批量获取所有反应物分子
  - `GetProductMolecules()` - 批量获取所有产物分子
  - `GetCatalystMolecules()` - 批量获取所有催化剂分子
- **Molecule 包增强**:
  - `NewMoleculeFromHandle()` - 从 Indigo 句柄创建分子对象

### 文档

- 新增反应分子访问文档 (docs/MOLECULE_ACCESS_IMPROVEMENTS.md)
- 新增反应分子操作示例 (examples/reaction/reaction_molecules.go)
- 完整的测试覆盖 (test/reaction/reaction_molecules_test.go)

### 改进

- 简化了反应中分子的访问方式
- 提供了批量操作接口提升效率
- 更好的文档和使用示例

## [0.4.3] - 2025-11-08

### 新增

- **格式支持增强**:
  - KET (Ketcher JSON) 格式支持
  - ChemAxon CXSmiles 格式增强
- **Molecule 包新增文件**:
  - molecule_atom.go: 20+ 原子/键操作方法
  - molecule_match.go: 模式匹配和高亮显示
- **Reaction 包辅助方法**:
  - reaction_helpers.go: 迭代器和布局辅助方法

### 改进

- 修复 SDF 保存实现
- 重组示例和测试结构
- 增强文档覆盖
- 改进 .gitignore 管理

### API 新增

- 原子和键操作方法
- 子结构匹配功能
- 反应迭代器辅助方法
- 分子和反应的 KET 格式方法

## [0.4.2] - 2025-11-08

### 新增

- **分子操作增强**:
  - 原子操作和键操作方法
  - 子结构匹配能力
- **格式转换增强**:
  - 扩展格式转换功能
  - 更多输出格式支持

### 改进

- 代码组织优化
- 示例代码重组
- 测试覆盖增强

## [0.4.0] - 2025-11-06

### 改进

- **构建系统优化**:
  - 统一跨平台 CGO 构建标志
  - 优化 macOS 链接器路径和包含目录
  - 改进第三方依赖结构组织
- **代码重构**:
  - 集中化 CGO 配置和对象创建逻辑
  - 简化反应模块 API（移除 Normalize、Standardize 和 Ionize 方法）
- **测试增强**:
  - 改进分子模块 base64 字符串验证测试

### 技术改进

- 更好的跨平台兼容性（Windows/Linux/macOS）
- 更清晰的 CGO 配置管理
- 更稳定的构建流程

## [0.3.0] - 2025-11-04

### 新增

- **Render 包**: 完整的分子和反应渲染功能
  - 支持 PNG、SVG、PDF 输出格式
  - 网格渲染支持
  - 丰富的渲染选项
  - 内存缓冲区渲染
- **文档系统重构**:
  - 新增英文 README (README_EN.md)
  - 新增贡献指南 (CONTRIBUTING.md)
  - 完整的文档中心 (docs/)
  - 模块级 README 文档
  - 详细的 API 参考
- **示例代码**:
  - example_render.go: 10个渲染示例
  - 完整的示例文档 (examples/README.md)

### 改进

- 统一了所有模块的 CGO 配置
- 支持多平台编译（Windows/Linux/macOS，x86_64/i386/aarch64）
- 优化了文档结构和导航

### 文档

- docs/API.md: 完整的 API 参考文档
- docs/INCHI.md: InChI 功能详细指南
- docs/SETUP.md: 环境配置完整指南
- render/README.md: 渲染功能文档

## [0.2.0] - 2025-11-03

### 新增

- **InChI 支持**: 使用 Indigo InChI 插件
  - InChI 生成
  - InChIKey 生成
  - 从 InChI 加载分子
  - 警告和日志信息
  - 辅助信息输出
- **Molecule 包完整重构**: 使用 CGO 绑定 Indigo 库
  - molecule.go: 核心分子结构
  - molecule_loader.go: 多格式加载
  - molecule_saver.go: 多格式保存
  - molecule_builder.go: 分子构建
  - molecule_properties.go: 属性计算
  - molecule_inchi.go: InChI 功能

### 改进

- 删除了旧的纯 Go 实现，统一使用 CGO
- 更好的错误处理和资源管理
- 使用 runtime.SetFinalizer 自动清理资源

### 修复

- 修复了分子属性计算中的 unsafe 导入问题
- 修复了元素常量缺失问题
- 修复了 InChI 模块初始化问题

### 文档

- molecule/README.md: 分子包完整文档
- examples/molecule/: 5个详细示例
- test/molecule/: 完整测试套件

## [0.1.0] - 2025-11-01

### 新增

- **Reaction 包**: 化学反应处理
  - 反应加载（SMILES, RXN 文件）
  - 反应保存（RXN, SMILES）
  - 自动原子映射（AAM）
  - 反应组件迭代
  - 反应标准化和归一化
- **CGO 集成**: 完整的 Indigo 库绑定
  - 跨平台支持
  - 自动内存管理
  - 错误处理封装

### 文档

- reaction/README.md: 反应包文档
- reaction/SETUP.md: 环境配置指南
- test/reaction/: 完整测试套件

## [0.0.1] - 2025-10-30

### 新增

- 项目初始化
- 基础项目结构
- Indigo 库集成准备
- go.mod 配置

---

## 版本说明

### 版本号格式

版本号格式为 `主版本号.次版本号.修订号`：

- **主版本号**: 不兼容的 API 修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

### 变更类型

- `新增` (Added): 新功能
- `改进` (Changed): 对现有功能的变更
- `弃用` (Deprecated): 即将移除的功能
- `移除` (Removed): 已移除的功能
- `修复` (Fixed): Bug 修复
- `安全` (Security): 安全相关修复

---

## 开发路线图

### v0.4.0 (计划中)

- [ ] 分子指纹计算
- [ ] Tanimoto 相似度
- [ ] 子结构匹配
- [ ] SMARTS 支持增强
- [ ] 性能优化

### v0.5.0 (计划中)

- [ ] 反应查询功能
- [ ] 反应相似度搜索
- [ ] 反应模板
- [ ] 批量处理工具

### v1.0.0 (目标)

- [ ] API 稳定
- [ ] 完整测试覆盖
- [ ] 性能基准测试
- [ ] 生产环境验证
- [ ] 完整文档

---

## 迁移指南

### 从 0.2.x 迁移到 0.3.x

**Render 包新增**:

```go
// 新功能：渲染分子
import "github.com/cx-luo/go-chem/render"

render.InitRenderer()
defer render.DisposeRenderer()

render.RenderToFile(mol.Handle(), "output.png")
```

**无破坏性变更**，所有 0.2.x 的代码继续工作。

### 从 0.1.x 迁移到 0.2.x

**重要变更**: Molecule 包完全重构。

**旧代码**（0.1.x，纯 Go）:

```go
// 不再支持
mol := molecule.NewMolecule()
```

**新代码**（0.2.x，CGO）:

```go
// 新的 API
mol, err := molecule.CreateMolecule()
if err != nil {
    log.Fatal(err)
}
defer mol.Close()  // 必须关闭
```

**主要变化**:

1. 所有函数返回 `error`
2. 必须显式 `Close()` 释放资源
3. 使用 `LoadMoleculeFromString()` 而不是构造函数
4. 新增 InChI 支持

---

## 贡献者

感谢所有为 go-chem 做出贡献的开发者！

- [@cx-luo](https://github.com/cx-luo) - 项目创建者和维护者

---

## 外部依赖

### Indigo Toolkit

- **版本**: 1.9.0+
- **许可证**: Apache License 2.0
- **网站**: <https://lifescience.opensource.epam.com/indigo/>

---

[Unreleased]: https://github.com/cx-luo/go-chem/compare/v0.4.4...HEAD
[0.4.4]: https://github.com/cx-luo/go-chem/compare/v0.4.3...v0.4.4
[0.4.3]: https://github.com/cx-luo/go-chem/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/cx-luo/go-chem/compare/v0.4.0...v0.4.2
[0.4.0]: https://github.com/cx-luo/go-chem/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/cx-luo/go-chem/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/cx-luo/go-chem/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/cx-luo/go-chem/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/cx-luo/go-chem/releases/tag/v0.0.1
