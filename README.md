# 规则引擎服务 (Rule Engine)

纯 Go 标准库实现的规则引擎后端服务，零第三方依赖。支持规则集、规则、条件和动作的编排与发布，输入求值（支持 9 种操作符）、版本管理、审计与执行记录。

## 运行说明

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，API Key 为 `default-ruleengine-key`（通过 Header `X-Api-Key` 或 Query `api_key` 传递）。

访问 `http://localhost:8080/api/health` 可检查服务状态。

## API 表格

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/rule-sets | 创建规则集 |
| GET | /api/rule-sets | 列出规则集（支持 name/status/keyword 筛选） |
| GET | /api/rule-sets/{id} | 获取规则集 |
| PUT | /api/rule-sets/{id} | 更新规则集 |
| DELETE | /api/rule-sets/{id} | 删除规则集 |
| POST | /api/rule-sets/{id}/publish | 发布（状态机 draft->published，生成版本快照+审计） |
| POST | /api/rule-sets/{id}/disable | 停用（状态机 published->disabled） |
| POST | /api/rules | 创建规则 |
| GET | /api/rules | 列出规则（支持 rule_set_id/status/enabled/keyword 筛选） |
| GET | /api/rules/{id} | 获取规则 |
| PUT | /api/rules/{id} | 更新规则 |
| DELETE | /api/rules/{id} | 删除规则 |
| POST | /api/rules/{id}/toggle | 启停规则（状态机 active->inactive） |
| POST | /api/conditions | 创建条件 |
| GET | /api/conditions | 列出条件（支持 field/operator/value_type/keyword 筛选） |
| GET | /api/conditions/{id} | 获取条件 |
| PUT | /api/conditions/{id} | 更新条件 |
| DELETE | /api/conditions/{id} | 删除条件 |
| POST | /api/actions | 创建动作 |
| GET | /api/actions | 列出动作（支持 action_type/target/keyword 筛选） |
| GET | /api/actions/{id} | 获取动作 |
| PUT | /api/actions/{id} | 更新动作 |
| DELETE | /api/actions/{id} | 删除动作 |
| GET | /api/evaluation-records | 列出评估记录（支持 rule_set_id/result/keyword 筛选） |
| GET | /api/evaluation-records/{id} | 获取评估记录 |
| DELETE | /api/evaluation-records/{id} | 删除评估记录 |
| POST | /api/rule-versions | 创建版本快照 |
| GET | /api/rule-versions | 列出版本（支持 rule_set_id/status/keyword 筛选） |
| GET | /api/rule-versions/{id} | 获取版本 |
| PUT | /api/rule-versions/{id} | 更新版本 |
| DELETE | /api/rule-versions/{id} | 删除版本 |
| POST | /api/rule-versions/{id}/rollback | 回滚版本 |
| GET | /api/execution-logs | 列出执行日志（支持 rule_set_id/rule_id/hit 筛选） |
| GET | /api/execution-logs/{id} | 获取执行日志 |
| DELETE | /api/execution-logs/{id} | 删除执行日志 |
| POST | /api/data-types | 创建数据类型 |
| GET | /api/data-types | 列出数据类型（支持 name/enabled/keyword 筛选） |
| GET | /api/data-types/{id} | 获取数据类型 |
| PUT | /api/data-types/{id} | 更新数据类型 |
| DELETE | /api/data-types/{id} | 删除数据类型 |
| POST | /api/audit-records | 创建审计记录 |
| GET | /api/audit-records | 列出审计记录（支持 operator/operation/target_type 筛选） |
| GET | /api/audit-records/{id} | 获取审计记录 |
| DELETE | /api/audit-records/{id} | 删除审计记录 |
| POST | /api/evaluate | 规则求值（输入 map 逐条件求值并执行动作） |
| POST | /api/batch/rules | 批量创建规则 |
| POST | /api/batch/conditions | 批量创建条件 |
| POST | /api/batch/actions | 批量创建动作 |
| POST | /api/batch/rules/delete | 批量删除规则 |
| POST | /api/cleanup/evaluation-records | 按时间清理评估记录 |
| POST | /api/cleanup/execution-logs | 按时间清理执行日志 |
| POST | /api/cleanup/audit-records | 按时间清理审计记录 |
| POST | /api/cleanup/rule-versions | 按保留数清理旧版本 |
| GET | /api/rule-sets/{id}/export | 导出规则集 JSON |
| POST | /api/rule-sets/import | 导入规则集 JSON |
| GET | /api/health | 健康检查 |
| GET | /api/health/dependencies | 依赖健康检查 |

## 实体列表

1. RuleSet（规则集）
2. Rule（规则）
3. Condition（条件）
4. Action（动作）
5. EvaluationRecord（评估记录）
6. RuleVersion（规则版本）
7. ExecutionLog（执行日志）
8. DataType（数据类型）
9. AuditRecord（审计记录）
10. ExecutionLog（执行记录）
