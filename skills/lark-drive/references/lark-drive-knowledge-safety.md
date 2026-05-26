# 知识资产整理安全策略

> 默认安全级别：只读或计划优先。任何会改变用户文档、目录或权限的操作，都必须先生成计划并让用户明确确认。

## 动作分级

| 分级 | 示例 | 策略 |
|-|-|-|
| Read-only | 列目录、查元数据、读 outline、读公开权限 | 可直接执行 |
| Plan-only | 生成整理计划、权限治理建议、报告草稿 | 可直接执行 |
| Confirmed write | 创建文件夹、创建知识库节点、移动文档、创建快捷方式、新建报告 | 必须先展示计划，用户确认后执行 |
| High-risk write | 删除、覆盖、改权限、转移 owner、移除成员、公开权限 patch | 第一版 workflow 不执行 |
| Unsupported | 单文档协作者全量枚举 | 明确说明当前能力无法验证 |

## 强制规则

- 盘点、整理建议、权限审计默认不修改任何资源。
- `organize-plan.json` 不是执行授权；只有用户明确说执行某个 plan，才进入执行阶段。
- 执行前必须展示 action 列表，包括 source、target、reason、risk、dry-run command。
- 执行写操作前先跑对应 `--dry-run`；dry-run 通过不等于用户已确认真实执行。
- 不自动执行 delete、overwrite、permission patch、member remove、owner transfer。
- 权限治理第一版只输出风险和建议，不自动收敛权限。
- 对无法确认的风险，写成 `unsupported_checks` 或 `requires_confirmation=true`，不要当作事实。

## 需要停下来问用户的情况

- 目标范围不明确，例如同时给出多个 folder/wiki URL 但未说明是否都处理。
- 计划包含跨空间移动、跨 Drive/Wiki 移动或大量移动。
- 目标目录已有同名节点/文件夹，无法判断是否复用。
- 用户要求“直接整理好”，但计划还没有展示和确认。
- 用户要求“治理权限”，但动作涉及改公开权限、移除成员、降权或转移 owner。
- 节点数量或文件数量超出上下文可处理范围，需要改为文件产物或缩小范围。

## 权限风险表述

表述必须区分：

- 事实：API 返回的字段，例如 `external_access=true`。
- 推断：基于规则得到的风险，例如“可能允许组织外传播”。
- 建议：下一步动作，例如“建议 owner 人工确认是否关闭外部访问”。
- 未覆盖：当前 API/CLI 不能验证的项，例如“无法枚举单文档显式协作者列表”。

示例：

```text
事实：link_share_entity=anyone_editable。
推断：互联网获得链接的人可能可编辑，属于高风险公开权限。
建议：请 owner 确认是否需要关闭链接分享或收敛为 tenant_readable。
未覆盖：未检查显式协作者列表，因为当前 CLI 没有 permission.members.list。
```
