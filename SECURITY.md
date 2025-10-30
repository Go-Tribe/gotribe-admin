# 安全指南

## 概述

GoTribe Admin 重视安全性，本文档提供了安全配置和最佳实践指南。

## 安全特性

### 1. 认证和授权
- **JWT Token 认证**: 使用 JWT 进行用户认证
- **RSA 加密**: 敏感数据传输使用 RSA 加密
- **密码加密**: 使用 BCrypt 进行密码哈希
- **RBAC 权限控制**: 基于角色的访问控制

### 2. 数据安全
- **SQL 注入防护**: 使用 GORM 参数化查询
- **XSS 防护**: 输入验证和输出编码
- **CSRF 防护**: 使用 CSRF Token
- **敏感数据加密**: 密码和敏感信息加密存储

### 3. 网络安全
- **HTTPS 支持**: 强制使用 HTTPS
- **CORS 配置**: 跨域请求控制
- **请求限流**: 防止暴力攻击
- **安全头**: 设置安全相关的 HTTP 头

## 安全配置

### 1. 生产环境配置

#### 数据库安全
```yaml
database:
  # 使用强密码
  password: "your_strong_password_here"
  # 启用 SSL
  sslmode: "require"
  # 限制连接数
  max-open-conns: 100
  max-idle-conns: 10
```

#### JWT 安全
```yaml
jwt:
  # 使用强密钥
  key: "your_very_strong_jwt_secret_key_here"
  # 设置合理的过期时间
  timeout: 2
  max-refresh: 24
```

#### 系统安全
```yaml
system:
  # 生产环境使用 release 模式
  mode: "release"
  # 禁用调试信息
  init-data: false
  # 启用数据迁移
  enable-migrate: true
```

### 2. 文件上传安全

#### 文件类型限制
```yaml
upload-file:
  # 限制文件类型
  allowed-exts: "jpg,jpeg,png,gif,pdf,doc,docx"
  # 限制文件大小
  max-size: 10485760  # 10MB
  # 启用病毒扫描
  enable-virus-scan: true
```

#### 文件存储安全
```yaml
upload-file:
  # 使用私有存储桶
  bucket: "your-private-bucket"
  # 设置访问权限
  acl: "private"
  # 启用加密
  server-side-encryption: true
```

### 3. 网络安全

#### HTTPS 配置
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 安全协议
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
}
```

#### 防火墙配置
```bash
# 只允许必要的端口
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw deny 3306/tcp  # MySQL (仅内网访问)
ufw deny 6379/tcp  # Redis (仅内网访问)
```

## 安全最佳实践

### 1. 密码策略
- 使用强密码（至少12位，包含大小写字母、数字、特殊字符）
- 定期更换密码
- 不要在多个服务中使用相同密码
- 使用密码管理器

### 2. 密钥管理
- 定期轮换 JWT 密钥
- 使用环境变量存储敏感信息
- 不要在代码中硬编码密钥
- 使用密钥管理服务

### 3. 数据库安全
- 使用最小权限原则
- 定期备份数据
- 启用数据库审计
- 使用连接加密

### 4. 应用安全
- 定期更新依赖包
- 使用安全扫描工具
- 实施代码审查
- 监控异常活动

### 5. 服务器安全
- 定期更新操作系统
- 使用非 root 用户运行应用
- 启用系统监控
- 实施入侵检测

## 安全监控

### 1. 日志监控
```yaml
logs:
  # 启用详细日志
  level: 0
  # 日志轮转
  max-size: 100
  max-backups: 30
  max-age: 7
  # 启用压缩
  compress: true
```

### 2. 异常检测
- 监控登录失败次数
- 监控异常请求
- 监控文件上传活动
- 监控数据库查询

### 3. 告警机制
- 设置安全事件告警
- 监控系统资源使用
- 监控网络流量
- 监控错误率

## 安全审计

### 1. 定期审计
- 每月进行安全审计
- 检查用户权限
- 审查访问日志
- 验证备份完整性

### 2. 漏洞扫描
- 使用自动化扫描工具
- 定期进行渗透测试
- 检查依赖包漏洞
- 验证安全配置

### 3. 合规性检查
- 检查数据保护合规性
- 验证访问控制
- 审查数据处理流程
- 确保日志完整性

## 应急响应

### 1. 安全事件处理
1. 立即隔离受影响系统
2. 收集证据和日志
3. 评估影响范围
4. 修复安全漏洞
5. 通知相关人员
6. 更新安全策略

### 2. 数据泄露处理
1. 立即停止数据泄露
2. 评估泄露影响
3. 通知相关用户
4. 报告监管机构
5. 加强安全措施
6. 进行事后分析

### 3. 系统恢复
1. 从备份恢复数据
2. 验证系统完整性
3. 更新安全补丁
4. 加强监控
5. 进行安全测试
6. 恢复正常服务

## 安全工具推荐

### 1. 漏洞扫描
- OWASP ZAP
- Nessus
- OpenVAS
- Nikto

### 2. 代码扫描
- SonarQube
- CodeQL
- Semgrep
- Gosec

### 3. 依赖扫描
- Snyk
- OWASP Dependency Check
- Trivy
- Grype

### 4. 运行时监控
- Falco
- OSSEC
- Wazuh
- ELK Stack

## 安全配置检查清单

### 部署前检查
- [ ] 修改默认密码
- [ ] 更新 JWT 密钥
- [ ] 配置 HTTPS
- [ ] 设置防火墙
- [ ] 启用日志记录
- [ ] 配置备份
- [ ] 设置监控
- [ ] 进行安全测试

### 运行时检查
- [ ] 监控异常活动
- [ ] 检查日志完整性
- [ ] 验证备份可用性
- [ ] 更新安全补丁
- [ ] 审查用户权限
- [ ] 检查系统资源
- [ ] 验证网络配置
- [ ] 测试应急响应

## 联系信息

如果您发现安全漏洞，请通过以下方式联系我们：

- 邮箱: security@gotribe.cn
- GitHub: [Security Advisories](https://github.com/go-tribe/gotribe-admin/security/advisories)

我们会在24小时内响应安全报告，并尽快修复漏洞。

## 免责声明

本安全指南仅供参考，不构成任何安全保证。用户需要根据自身环境和需求进行适当的安全配置和测试。
