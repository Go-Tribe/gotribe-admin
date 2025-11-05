# 项目文档索引

本页汇总项目文档入口，区分人写文档与自动生成内容，方便查阅与维护。

## 使用指南（Guides）
- 架构概览：`guides/ARCHITECTURE.md`
- 中文总览：`guides/README_CN.md`
- 国际化示例：`guides/internationalization_example.md`

## 参考资料（Reference）
- API 说明：`reference/API.md`
- 版本变更（Changelog）：`reference/CHANGELOG.md`

## Swagger 生成产物（Auto-generated）
- OpenAPI YAML：`swagger/swagger.yaml`
- OpenAPI JSON：`swagger/swagger.json`
- Go 文档包：`swagger/docs.go`（在 `gotribe-admin.go` 通过 `_ "gotribe-admin/docs/swagger"` 旁路导入）

## 资源
- 截图资源：`images/`

## 目录约定
- 人写文档存放在：`docs/guides/` 与 `docs/reference/`
- 自动生成内容存放在：`docs/swagger/`

如需进一步优化结构或添加更多索引，请在此页补充相应入口。
