# 知识资产整理 Workflow 总览

> **前置条件：** 先阅读 [`../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) 了解认证、全局参数和安全规则。

## 适用场景

- 盘点某个云空间文件夹、知识库、知识库节点或文档库下的文档和目录。
- 整理散乱文档：归类、生成目录结构、移动计划、快捷方式计划、导航页建议。
- 治理知识资产：重复、空标题、孤立、未归档、命名混乱、内容缺失或过期。
- 审计权限风险：Wiki 空间成员、文档公开链接、组织外访问、分享/复制/下载设置。

## 目标范围模型

| 用户目标 | 视作 | 主入口 |
|-|-|-|
| 云空间文件夹 | Drive folder tree | `drive files list` / `drive +search` |
| 知识库 / 知识库节点 | Wiki space/node tree | `wiki +node-list` / `wiki +node-get` |
| 文档库 / my_library | Wiki personal library | `wiki +node-list --space-id my_library --as user` |
| 文档正文或标题结构 | Docs content | `docs +fetch --api-version v2` |
| 报告或台账 | Docs / Sheets | `docs +create` / `sheets +create` |

不要把“文档库”当成 Drive 根目录；它应走 Wiki personal library。

## Recipe 路由

| 用户意图 | 读取文件 | 产物 |
|-|-|-|
| 盘点、梳理、导出清单、看目录结构 | [`lark-drive-knowledge-inventory.md`](lark-drive-knowledge-inventory.md) | `inventory.json` |
| 整理、归类、组织目录、生成移动计划 | [`lark-drive-knowledge-organize.md`](lark-drive-knowledge-organize.md) | `organize-plan.json` |
| 权限风险、公开链接、组织外访问、成员过多 | [`lark-drive-knowledge-permission-audit.md`](lark-drive-knowledge-permission-audit.md) | `permission-audit.json` |

复合意图按链路执行：

```text
盘点并整理 -> inventory -> organize
盘点并审计权限 -> inventory -> permission-audit
整理并输出报告 -> inventory -> organize -> report artifact
```

## 执行原则

- 大范围结果默认写入本地 artifact，不把完整清单塞进聊天。
- 每次 run 使用独立目录：`./lark-drive-knowledge/<run-id>/`。
- 工作流之间通过 artifact 串联，优先复用已有 `inventory.json`，不要无故重复爬取。
- 所有判断必须区分事实、推断、建议；没有证据的内容写入 `warnings` 或 `unsupported_checks`。
- 原生 API 调用前必须先运行 `lark-cli schema <service>.<resource>.<method>` 校验参数结构。
- 写操作和权限治理必须遵守 [`lark-drive-knowledge-safety.md`](lark-drive-knowledge-safety.md)。

## 当前能力边界

| 能力 | 状态 | 说明 |
|-|-|-|
| 知识库节点盘点 | 支持 | 递归 `wiki +node-list` |
| 云空间文件夹盘点 | 支持 | 递归 `drive files list` |
| 文档标题结构读取 | 支持 | `docs +fetch --scope outline` |
| Wiki 空间成员审计 | 支持 | `wiki +member-list` |
| 文档公开权限审计 | 支持 | `drive permission.public get` |
| 单文档协作者全量枚举 | 暂不支持 | 当前 Drive permission members 只有 `auth/create/transfer_owner` |
| 自动删除、覆盖、降权、改权限 | 第一版不执行 | 只输出计划和建议 |
