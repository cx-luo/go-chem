# 贡献指南 / Contributing Guide

[English](#english) | [简体中文](#简体中文)

---

## English

Thank you for your interest in contributing to go-chem! This document provides guidelines for contributing to the project.

### How to Contribute

#### Reporting Bugs

If you find a bug, please create an issue with:

1. **Clear title**: Briefly describe the problem
2. **Environment**: OS, Go version, architecture
3. **Steps to reproduce**: Detailed steps to reproduce the issue
4. **Expected behavior**: What should happen
5. **Actual behavior**: What actually happens
6. **Code sample**: Minimal code to reproduce the issue

Example:

```markdown
**Title**: Memory leak when cloning molecules

**Environment**:
- OS: Windows 10
- Go: 1.21.0
- Architecture: amd64

**Steps to reproduce**:
1. Load a molecule from SMILES
2. Clone it 1000 times
3. Observe memory usage

**Expected**: Memory should be freed after Close()
**Actual**: Memory keeps increasing

**Code**:
\`\`\`go
for i := 0; i < 1000; i++ {
    clone, _ := mol.Clone()
    clone.Close()
}
\`\`\`
```

#### Suggesting Features

For feature requests, please:

1. **Check existing issues**: Avoid duplicates
2. **Describe use case**: Why is this feature needed?
3. **Provide examples**: How would it be used?
4. **Consider alternatives**: Are there other ways to achieve this?

#### Pull Requests

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/my-feature`
3. **Make your changes**: Follow coding standards
4. **Add tests**: Ensure your code is tested
5. **Update documentation**: If needed
6. **Commit with clear messages**: Describe what and why
7. **Push to your fork**: `git push origin feature/my-feature`
8. **Open a Pull Request**: Provide detailed description

### Coding Standards

#### Go Code Style

Follow standard Go conventions:

```go
// ✅ Good
func LoadMoleculeFromString(smiles string) (*Molecule, error) {
    if smiles == "" {
        return nil, fmt.Errorf("empty SMILES string")
    }
    // Implementation
}

// ❌ Bad
func loadMol(s string) *Molecule {
    // No error handling
}
```

#### Documentation

- All exported functions must have comments
- Comments should start with the function name
- Provide usage examples for complex functions

```go
// LoadMoleculeFromString loads a molecule from a SMILES string.
// It returns an error if the SMILES string is invalid or empty.
//
// Example:
//   mol, err := LoadMoleculeFromString("CCO")
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer mol.Close()
func LoadMoleculeFromString(smiles string) (*Molecule, error) {
    // Implementation
}
```

#### Testing

- Write tests for new functionality
- Aim for high test coverage
- Use table-driven tests when appropriate

```go
func TestLoadMoleculeFromString(t *testing.T) {
    tests := []struct {
        name    string
        smiles  string
        wantErr bool
    }{
        {"valid ethanol", "CCO", false},
        {"valid benzene", "c1ccccc1", false},
        {"empty string", "", true},
        {"invalid smiles", "C(", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mol, err := LoadMoleculeFromString(tt.smiles)
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadMoleculeFromString() error = %v, wantErr %v", err, tt.wantErr)
            }
            if mol != nil {
                mol.Close()
            }
        })
    }
}
```

#### Error Handling

- Always check and handle errors
- Provide descriptive error messages
- Use `fmt.Errorf` with context

```go
// ✅ Good
if err != nil {
    return nil, fmt.Errorf("failed to load molecule from string: %w", err)
}

// ❌ Bad
if err != nil {
    return nil, err  // Lost context
}
```

#### Resource Management

- Always close resources
- Use `defer` for cleanup
- Add finalizers as safety net

```go
// ✅ Good
func ProcessMolecule(smiles string) error {
    mol, err := LoadMoleculeFromString(smiles)
    if err != nil {
        return err
    }
    defer mol.Close()  // Ensures cleanup

    // Process molecule
    return nil
}

// ❌ Bad
func ProcessMolecule(smiles string) error {
    mol, _ := LoadMoleculeFromString(smiles)
    // Forgot to close, potential memory leak
    return nil
}
```

### Development Workflow

1. **Setup environment**: Follow [SETUP.md](docs/SETUP.md)
2. **Create feature branch**: `git checkout -b feature/my-feature`
3. **Make changes**: Implement your feature
4. **Run tests**: `go test ./...`
5. **Run linter**: `golangci-lint run`
6. **Update docs**: If API changed
7. **Commit**: Use clear commit messages
8. **Push**: `git push origin feature/my-feature`
9. **Create PR**: Open pull request on GitHub

### Commit Message Guidelines

Use conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

Examples:

```
feat(molecule): add InChI generation support

- Add ToInChI method
- Add ToInChIKey method
- Update tests

Closes #42
```

```
fix(render): correct PNG output buffer handling

The buffer was not properly initialized, causing crashes
on some platforms.

Fixes #56
```

### Code Review Process

1. **Automated checks**: CI must pass
2. **Code review**: Maintainer reviews code
3. **Feedback**: Address review comments
4. **Approval**: Maintainer approves
5. **Merge**: PR is merged

### Getting Help

- **Documentation**: Check [docs/](docs/)
- **Issues**: Search existing issues
- **Discussions**: Start a discussion on GitHub
- **Email**: <chengxiang.luo@foxmail.com>

---

## 简体中文

感谢您对 go-chem 项目的关注！本文档提供了贡献指南。

### 如何贡献

#### 报告 Bug

如果发现 bug，请创建 issue 并包含：

1. **清晰的标题**: 简要描述问题
2. **环境信息**: 操作系统、Go 版本、架构
3. **重现步骤**: 详细的重现步骤
4. **期望行为**: 应该发生什么
5. **实际行为**: 实际发生了什么
6. **代码示例**: 可重现问题的最小代码

#### 功能建议

对于功能请求，请：

1. **检查现有 issue**: 避免重复
2. **描述用例**: 为什么需要这个功能？
3. **提供示例**: 如何使用？
4. **考虑替代方案**: 是否有其他方法实现？

#### Pull Request

1. **Fork 仓库**
2. **创建功能分支**: `git checkout -b feature/my-feature`
3. **进行修改**: 遵循编码规范
4. **添加测试**: 确保代码有测试
5. **更新文档**: 如有需要
6. **清晰的提交信息**: 描述做了什么和为什么
7. **推送到 fork**: `git push origin feature/my-feature`
8. **开启 Pull Request**: 提供详细描述

### 编码规范

#### Go 代码风格

遵循标准 Go 规范：

```go
// ✅ 好的做法
func LoadMoleculeFromString(smiles string) (*Molecule, error) {
    if smiles == "" {
        return nil, fmt.Errorf("空的 SMILES 字符串")
    }
    // 实现
}

// ❌ 不好的做法
func loadMol(s string) *Molecule {
    // 没有错误处理
}
```

#### 文档注释

- 所有导出函数必须有注释
- 注释应以函数名开头
- 为复杂函数提供使用示例

#### 测试

- 为新功能编写测试
- 追求高测试覆盖率
- 适当使用表驱动测试

#### 错误处理

- 总是检查和处理错误
- 提供描述性错误信息
- 使用 `fmt.Errorf` 添加上下文

#### 资源管理

- 总是关闭资源
- 使用 `defer` 进行清理
- 添加 finalizer 作为安全网

### 开发流程

1. **配置环境**: 参考 [SETUP.md](docs/SETUP.md)
2. **创建功能分支**: `git checkout -b feature/my-feature`
3. **进行修改**: 实现功能
4. **运行测试**: `go test ./...`
5. **运行 linter**: `golangci-lint run`
6. **更新文档**: 如果 API 改变
7. **提交**: 使用清晰的提交信息
8. **推送**: `git push origin feature/my-feature`
9. **创建 PR**: 在 GitHub 上开启 pull request

### 提交信息规范

使用约定式提交格式：

```
<类型>(<范围>): <主题>

<正文>

<页脚>
```

类型：

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档变更
- `style`: 代码格式（不影响代码逻辑）
- `refactor`: 代码重构
- `test`: 添加测试
- `chore`: 维护任务

示例：

```
feat(molecule): 添加 InChI 生成支持

- 添加 ToInChI 方法
- 添加 ToInChIKey 方法
- 更新测试

Closes #42
```

### 代码审查流程

1. **自动检查**: CI 必须通过
2. **代码审查**: 维护者审查代码
3. **反馈**: 处理审查意见
4. **批准**: 维护者批准
5. **合并**: PR 被合并

### 获取帮助

- **文档**: 查看 [docs/](docs/)
- **Issues**: 搜索现有 issue
- **讨论**: 在 GitHub 上开启讨论
- **邮件**: <chengxiang.luo@foxmail.com>

---

## 行为准则 / Code of Conduct

### 我们的承诺 / Our Pledge

为了营造一个开放和友好的环境，我们承诺让每个人都能参与我们的项目和社区，不论其年龄、体型、残疾、种族、性别认同和表达、经验水平、教育程度、社会经济地位、国籍、个人外貌、种族、宗教或性认同和性取向。

### 我们的标准 / Our Standards

积极行为的例子：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表现出同理心

不可接受的行为：

- 使用性化的语言或图像
- 挑衅、侮辱或贬损的评论
- 公开或私下的骚扰
- 未经许可发布他人的私人信息
- 其他在专业场合可被认为不适当的行为

### 执行 / Enforcement

如有任何不当行为，请联系项目维护者：<chengxiang.luo@foxmail.com>

---

Thank you for contributing to go-chem! 🙏

感谢您为 go-chem 做出贡献！🙏
