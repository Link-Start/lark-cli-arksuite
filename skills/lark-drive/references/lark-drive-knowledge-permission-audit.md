# 权限风险审计 Workflow

> 前置条件：先读 [`lark-drive-knowledge-overview.md`](lark-drive-knowledge-overview.md)、[`lark-drive-knowledge-artifacts.md`](lark-drive-knowledge-artifacts.md) 和 [`lark-drive-knowledge-safety.md`](lark-drive-knowledge-safety.md)。第一版只审计和建议，不自动改权限。

## 目标

审计 云空间 / 知识库 / 文档库范围内的权限风险，重点覆盖：

- Wiki 空间成员和角色。
- 文档公开权限、外部访问、链接分享、复制/下载/打印限制。
- 当前能力无法验证的权限项。

## 当前能力边界

| 检查项 | 状态 | 命令 |
|-|-|-|
| Wiki 空间成员 | 支持 | `wiki +member-list` |
| 文档公开权限 | 支持 | `drive permission.public get`，仅限 schema 支持的文档类型 |
| Drive 文件夹公开权限 | 暂不支持 | `drive.permission.public.get` 不支持 `type=folder` |
| 当前用户/应用是否具备某权限 | 支持 | `drive permission.members auth` |
| 单文档显式协作者全量枚举 | 暂不支持 | 当前无 `drive.permission.members.list` |
| 自动关闭公开权限或降权 | 第一版不执行 | 只输出建议 |

## Step 1: 输入范围

优先消费 `inventory.json`。如果用户只给一个 URL，先按 [`lark-drive-knowledge-inventory.md`](lark-drive-knowledge-inventory.md) 做最小盘点。

## Step 2: Wiki 空间成员审计

涉及知识库空间或文档库 / `my_library` 时：

```bash
lark-cli wiki +member-list \
  --space-id "<space_id_or_my_library>" \
  --page-all --page-limit 0 --format json --as user
```

审计信号：

- admin 数量异常多。
- 成员包含部门、群或开放范围较大的 member_type。
- 文档库成员结果缺失或不适用时写入 `warnings`，不要推断。

## Step 3: 文档公开权限审计

先查看 schema：

```bash
lark-cli schema drive.permission.public.get --format json
```

只对 inventory 中 `type` 属于以下集合的条目查询：

```text
doc, docx, sheet, bitable, file, wiki, mindnote, minutes, slides
```

`type=folder` 的 Drive 文件夹不要调用 `drive.permission.public get`。将该 item 写入 `unsupported_checks` 或 `warnings`，例如 `drive_folder_public_permission_unsupported`，并继续审计其他条目。

```bash
lark-cli drive permission.public get \
  --params '{"token":"<token>","type":"<type>"}' \
  --format json --as user
```

Token/type 选择：

- 云空间文件：用 Drive `token` 和 `type=file`。
- 云空间文件夹：不执行 `permission.public.get`，记录为未覆盖检查。
- 知识库节点权限：优先用 `node_token` 和 `type=wiki`。
- 底层文档权限：用 `obj_token` 和 `obj_type`。

如果某类 token 查询失败，不要中断全局审计；写入该 item 的 `warnings`。如果是 `type=folder`，不要把它当作 API 失败，而是写入 `unsupported_checks`。

## 风险规则

| 字段 | 高风险 | 中风险 |
|-|-|-|
| `link_share_entity` | `anyone_editable`, `anyone_readable` | `tenant_editable` |
| `external_access` | `true` 且链接或分享范围较宽 | `true` |
| `share_entity` | `anyone` | `same_tenant` |
| `security_entity` | `anyone_can_view` | `anyone_can_edit` |
| Wiki member role | admin 过多 | 部门/群成员需确认 |

每条风险都必须写清楚：

- `fact`：API 返回字段。
- `inference`：为什么可能有风险。
- `suggestion`：建议 owner 或管理员确认的动作。
- `evidence`：token、URL、path、字段名和值。
- `requires_confirmation=true`。

## 输出

写入 `permission-audit.json`，聊天回复只给摘要：

- 审计对象数量。
- 高/中/低风险数量。
- 公开权限风险 Top N。
- Wiki 成员风险摘要。
- 未覆盖检查：必须包含 `explicit_collaborator_list`，说明当前 CLI 无法枚举单文档显式协作者列表。
- 如果 inventory 中包含 `type=folder`，未覆盖检查必须包含 `drive_folder_public_permission`，说明 `drive.permission.public.get` 不支持 Drive 文件夹。

## 禁止事项

- 不自动调用 `permission.public.patch`。
- 不自动调用 `permission.members.create`、`transfer_owner` 或任何成员删除/降权接口。
- 不把“无法枚举协作者”说成“没有协作者风险”。
- 不把模型推断写成事实。
