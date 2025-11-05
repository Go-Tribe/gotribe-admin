# 贡献指南

感谢您对 GoTribe Admin 项目的关注！我们欢迎任何形式的贡献，包括但不限于：

- 报告 Bug
- 提出新功能建议
- 提交代码修复
- 改进文档
- 分享使用经验

## 开发环境搭建

### 前置要求

- Go 1.20 或更高版本
- MySQL 8.0+ 或 PostgreSQL 13+
- Redis 6.0+
- Node.js 18+ (用于前端开发)
- Git

### 环境配置

1. **Fork 并克隆项目**

```bash
# Fork 项目到您的 GitHub 账户
# 然后克隆您的 Fork
git clone https://github.com/your-username/gotribe-admin.git
cd gotribe-admin
```

2. **安装依赖**

```bash
# 安装 Go 依赖
go mod tidy

# 安装前端依赖 (如果需要修改前端)
cd web/admin
npm install
```

3. **配置数据库**

```bash
# 复制配置文件
cp config.tmp.yml config.yml

# 编辑配置文件，设置数据库连接信息
vim config.yml
```

4. **初始化数据库**

```bash
# 运行数据库迁移
make migrate

# 初始化基础数据
make seed
```

5. **启动项目**

```bash
# 启动后端服务
make run

# 启动前端服务 (如果需要)
cd web/admin
npm run dev
```

## 开发规范

### 代码规范

1. **Go 代码规范**
   - 遵循 [Go 官方代码规范](https://golang.org/doc/effective_go.html)
   - 使用 `gofmt` 格式化代码
   - 使用 `golint` 检查代码质量
   - 添加必要的注释

2. **命名规范**
   - 包名：小写，简短，有意义
   - 函数名：驼峰命名，公开函数首字母大写
   - 变量名：驼峰命名，简洁明了
   - 常量名：全大写，下划线分隔

3. **注释规范**
   - 公开的函数、类型、变量必须有注释
   - 注释应该说明功能，而不是实现
   - 使用 `//` 进行单行注释
   - 使用 `/* */` 进行多行注释

### 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型说明：**
- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例：**
```bash
git commit -m "feat(user): add user avatar upload functionality"
git commit -m "fix(api): resolve user login validation issue"
git commit -m "docs: update API documentation"
```

### 分支管理

1. **主分支**
   - `main`: 主分支，用于发布稳定版本
   - `develop`: 开发分支，用于集成新功能

2. **功能分支**
   - `feature/功能名称`: 新功能开发
   - `bugfix/问题描述`: Bug 修复
   - `hotfix/问题描述`: 紧急修复

3. **分支命名规范**
   ```bash
   # 新功能
   feature/user-management
   feature/payment-integration

   # Bug 修复
   bugfix/login-validation
   bugfix/file-upload-error

   # 紧急修复
   hotfix/security-patch
   hotfix/critical-bug
   ```

## 贡献流程

### 1. 报告问题

如果您发现了 Bug 或有功能建议，请：

1. 检查 [Issues](https://github.com/go-tribe/gotribe-admin/issues) 是否已存在
2. 创建新的 Issue，详细描述问题或建议
3. 使用合适的标签标记 Issue

**Issue 模板：**
```markdown
## 问题描述
简要描述问题或建议

## 重现步骤
1. 进入 '...'
2. 点击 '...'
3. 滚动到 '...'
4. 看到错误

## 预期行为
描述您期望的行为

## 实际行为
描述实际发生的行为

## 环境信息
- OS: [e.g. macOS, Windows, Linux]
- Go Version: [e.g. 1.20.0]
- Database: [e.g. MySQL 8.0, PostgreSQL 13]
- Browser: [e.g. Chrome, Firefox, Safari]

## 附加信息
添加任何其他相关信息
```

### 2. 提交代码

1. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **开发功能**
   - 编写代码
   - 添加测试
   - 更新文档
   - 确保代码通过所有检查

3. **提交代码**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

4. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 详细描述您的更改
   - 关联相关的 Issue

### 3. 代码审查

所有提交的代码都会经过审查：

1. **自动化检查**
   - 代码格式检查
   - 单元测试
   - 集成测试
   - 安全扫描

2. **人工审查**
   - 代码质量
   - 功能正确性
   - 性能影响
   - 安全性

3. **审查反馈**
   - 审查者会提供反馈
   - 需要根据反馈进行修改
   - 修改后重新提交

## 测试指南

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/app/controller

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 集成测试

```bash
# 运行集成测试
go test -tags=integration ./...

# 运行特定测试
go test -run TestUserController ./internal/app/controller
```

### 性能测试

```bash
# 运行性能测试
go test -bench=. ./...

# 运行特定性能测试
go test -bench=BenchmarkUserController ./internal/app/controller
```

## 文档贡献

### 代码文档

- 为公开的 API 添加注释
- 使用 `godoc` 生成文档
- 确保注释准确、简洁

### 用户文档

- 更新 README.md
- 完善 API 文档
- 添加使用示例
- 更新架构文档

### 文档规范

- 使用 Markdown 格式
- 保持文档结构清晰
- 添加必要的示例代码
- 定期更新过时信息

## 发布流程

### 版本号规范

我们使用 [Semantic Versioning](https://semver.org/)：

- `MAJOR`: 不兼容的 API 修改
- `MINOR`: 向下兼容的功能性新增
- `PATCH`: 向下兼容的问题修正

### 发布步骤

1. **准备发布**
   ```bash
   # 更新版本号
   # 更新 CHANGELOG.md
   # 更新文档
   ```

2. **创建发布**
   ```bash
   # 创建标签
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **发布说明**
   - 列出新功能
   - 说明 Bug 修复
   - 提供升级指南
   - 标记破坏性变更

## 社区参与

### 讨论

- 使用 [GitHub Discussions](https://github.com/go-tribe/gotribe-admin/discussions) 进行讨论
- 参与技术交流
- 分享使用经验
- 提出改进建议

### 支持

- 帮助其他用户解决问题
- 回答技术问题
- 提供使用建议
- 分享最佳实践

## 行为准则

我们致力于为每个人提供友好、安全、包容的环境。请遵守以下准则：

1. **尊重他人**
   - 使用友好和包容的语言
   - 尊重不同的观点和经验
   - 接受建设性的批评

2. **专业行为**
   - 专注于对社区最有利的事情
   - 避免个人攻击
   - 保持专业和礼貌

3. **包容性**
   - 欢迎不同背景的贡献者
   - 尊重不同的技能水平
   - 提供帮助和支持

## 许可证

通过贡献代码，您同意您的贡献将在 MIT 许可证下发布。

## 联系方式

如果您有任何问题或建议，请通过以下方式联系我们：

- GitHub Issues: [https://github.com/go-tribe/gotribe-admin/issues](https://github.com/go-tribe/gotribe-admin/issues)
- GitHub Discussions: [https://github.com/go-tribe/gotribe-admin/discussions](https://github.com/go-tribe/gotribe-admin/discussions)
- Email: [my@dengmengmian.com](mailto:my@dengmengmian.com)

感谢您的贡献！
