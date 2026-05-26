# 知识资产盘点 Workflow

> 前置条件：先读 [`lark-drive-knowledge-overview.md`](lark-drive-knowledge-overview.md) 和 [`lark-drive-knowledge-artifacts.md`](lark-drive-knowledge-artifacts.md)。涉及 Wiki 或文档库时再读 [`../../lark-wiki/SKILL.md`](../../lark-wiki/SKILL.md)。

## 目标

对云空间文件夹、知识库、知识库子树或文档库做结构化盘点，生成 `inventory.json`，作为整理、权限审计和报告的事实底座。

## Step 1: 归一化范围

- 云空间文件夹 URL：提取 `folder_token`，或用 `drive +inspect` 确认类型。
- 知识库 URL：用 `wiki +node-get` 或 `drive +inspect` 获取 `space_id`、`node_token`、`obj_type`、`obj_token`。
- `my_library` / 文档库：固定走 `wiki +node-list --space-id my_library --as user`。

记录到 `scope.json`。

## Step 2: 盘点云空间文件夹

先查看 schema：

```bash
lark-cli schema drive.files.list --format json
```

读取直接子项：

```bash
lark-cli drive files list \
  --params '{"folder_token":"<folder_token>","page_size":200}' \
  --page-all --page-limit 0 --format json --as user
```

对返回的 `type=folder` 子项递归调用同一命令。记录字段：

- `name` -> `title`
- `token`
- `type`
- `url`
- `parent_token`
- `owner_id`
- `created_time`
- `modified_time`

如需按关键词或时间补充范围，可用 `drive +search`，但目录树以 `drive files list` 为准。

## Step 3: 盘点知识库或文档库

读取根层：

```bash
lark-cli wiki +node-list \
  --space-id "<space_id_or_my_library>" \
  --page-all --page-limit 0 --format json --as user
```

对 `has_child=true` 的节点递归：

```bash
lark-cli wiki +node-list \
  --space-id "<space_id_or_my_library>" \
  --parent-node-token "<node_token>" \
  --page-all --page-limit 0 --format json --as user
```

必要时补节点详情：

```bash
lark-cli wiki +node-get --node-token "<node_token>" --format json --as user
```

记录字段：

- `title`
- `node_token`
- `obj_token`
- `obj_type`
- `node_type`
- `parent_node_token`
- `has_child`
- `space_id`
- `owner`（如果详情返回）

## Step 4: 可选读取内容结构

只有用户需要“按内容归类”“识别主题”“生成导航页”时才读取正文结构。优先读 outline，不全量读正文：

```bash
lark-cli docs +fetch \
  --api-version v2 \
  --doc "<url_or_token>" \
  --scope outline \
  --max-depth 3 \
  --doc-format markdown \
  --format json --as user
```

失败时保留结构盘点结果，并在 `warnings` 中标记 `content_outline_failed`。

## Step 5: 输出 inventory.json

把所有条目写入 `./lark-drive-knowledge/<run-id>/inventory.json`。聊天回复只输出摘要：

- 总数
- 按 source/type 统计
- 空标题数量
- 重复标题候选数量
- 未读取或失败的节点数量
- artifact 路径

## 停止条件

- 没有权限读取根节点或根文件夹：停止并给出授权建议。
- 分页失败或 cursor 不前进：保留已读结果，写入 `warnings`。
- 结果规模过大：停止深挖正文，只输出结构清单并提示用户缩小范围。
- 知识库和云空间混合范围不清楚：先让用户确认是否都处理。
