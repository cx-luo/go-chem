# go-indigo 文档中心

欢迎来到 go-indigo 的文档中心。这里包含了所有详细的技术文档和使用指南。

## 📖 文档索引

### 核心模块文档

1. **[API 参考](API.md)** - 完整的 API 文档
2. **[InChI 功能](INCHI.md)** - InChI 和 InChIKey 完整指南
3. **[渲染指南](RENDERING.md)** - 分子和反应结构渲染
4. **[环境配置](SETUP.md)** - CGO 环境设置详解

### 快速入门

- **[5分钟快速开始](QUICKSTART.md)** - 最快上手指南
- **[常见问题 FAQ](FAQ.md)** - 常见问题解答
- **[最佳实践](BEST_PRACTICES.md)** - 推荐的使用模式

### 高级主题

- **[性能优化](PERFORMANCE.md)** - 性能调优指南
- **[CGO 深入](CGO_ADVANCED.md)** - CGO 高级用法
- **[扩展开发](EXTENSIONS.md)** - 如何扩展功能

## 📚 按功能分类

### 分子处理

| 文档 | 描述 |
|------|------|
| [分子基础](molecule/BASICS.md) | 分子结构基础概念 |
| [分子加载](molecule/LOADING.md) | 各种格式的加载方法 |
| [分子构建](molecule/BUILDING.md) | 从头构建分子 |
| [分子属性](molecule/PROPERTIES.md) | 属性计算详解 |

### 反应处理

| 文档 | 描述 |
|------|------|
| [反应基础](reaction/BASICS.md) | 化学反应基础 |
| [原子映射](reaction/AAM.md) | 自动原子映射 |
| [反应查询](reaction/QUERY.md) | 反应查询和搜索 |

### 渲染和可视化

| 文档 | 描述 |
|------|------|
| [渲染基础](render/BASICS.md) | 渲染功能概述 |
| [自定义样式](render/STYLING.md) | 自定义渲染样式 |
| [批量渲染](render/BATCH.md) | 批量渲染技巧 |

## 🎯 学习路径

### 初学者路径

1. 阅读 [快速开始](QUICKSTART.md)
2. 了解 [环境配置](SETUP.md)
3. 查看 [molecule 基础示例](../examples/molecule/)
4. 尝试 [常见用例](USE_CASES.md)

### 进阶路径

1. 深入 [API 参考](API.md)
2. 学习 [性能优化](PERFORMANCE.md)
3. 理解 [CGO 机制](CGO_ADVANCED.md)
4. 探索 [扩展开发](EXTENSIONS.md)

## 📝 文档规范

### 文档结构

每个文档应包含：

- **标题和简介**: 清晰说明文档目的
- **目录**: 便于快速导航
- **代码示例**: 实际可运行的代码
- **参数说明**: 详细的参数描述
- **错误处理**: 常见错误和解决方案
- **相关链接**: 关联文档的链接

### 代码示例规范

```go
// ✅ 好的示例：包含错误处理和资源释放
func GoodExample() {
    mol, err := molecule.LoadMoleculeFromString("CCO")
    if err != nil {
        log.Fatalf("加载失败: %v", err)
    }
    defer mol.Close()

    // 使用分子
    smiles, _ := mol.ToSmiles()
    fmt.Println(smiles)
}

// ❌ 不好的示例：忽略错误，不释放资源
func BadExample() {
    mol, _ := molecule.LoadMoleculeFromString("CCO")
    smiles, _ := mol.ToSmiles()
    fmt.Println(smiles)
}
```

## 🔍 文档搜索

如果你在寻找特定功能，可以：

1. 使用 GitHub 仓库搜索功能
2. 查看 [API 索引](API.md#索引)
3. 浏览 [示例代码](../examples/)
4. 查阅 [FAQ](FAQ.md)

## 📚 外部资源

### Indigo 官方文档

- [Indigo Toolkit](https://lifescience.opensource.epam.com/indigo/)
- [Indigo API Reference](https://lifescience.opensource.epam.com/indigo/api/)
- [Indigo GitHub](https://github.com/epam/Indigo)

### 化学信息学资源

- [SMILES 规范](https://www.daylight.com/dayhtml/doc/theory/theory.smiles.html)
- [InChI FAQ](https://www.inchi-trust.org/inchi-faq/)
- [MDL MOL 格式](http://c4.cabrillo.edu/404/ctfile.pdf)

### Go 语言资源

- [Go 官方文档](https://golang.org/doc/)
- [CGO 文档](https://golang.org/cmd/cgo/)
- [Effective Go](https://golang.org/doc/effective_go)

## 🤝 贡献文档

我们欢迎文档贡献！如果你发现：

- 文档中的错误或不准确之处
- 缺失的功能说明
- 需要更多示例的地方
- 可以改进的表述

请提交 Issue 或 Pull Request。

### 文档贡献指南

1. Fork 仓库
2. 在 `docs/` 目录下创建或修改文档
3. 使用 Markdown 格式
4. 包含代码示例
5. 提交 Pull Request

## 📞 获取帮助

如果文档没有解答你的问题：

1. 查看 [FAQ](FAQ.md)
2. 在 GitHub 创建 Issue
3. 发送邮件至 <chengxiang.luo@foxmail.com>

## 📊 文档更新日志

- **2025-11-04**: 创建文档中心，重组文档结构
- **2025-11-03**: 添加 InChI 相关文档
- **2025-11-02**: 添加 Render 模块文档
- **2025-11-01**: 初始化文档

---

💡 **提示**: 建议按照学习路径循序渐进地阅读文档，这样可以更系统地掌握 go-indigo 的使用。
