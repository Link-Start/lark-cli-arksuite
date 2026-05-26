# 散乱知识整理 Workflow

> 前置条件：先有 `inventory.json`，并阅读 [`lark-drive-knowledge-safety.md`](lark-drive-knowledge-safety.md)。第一版默认只生成整理计划，不直接执行。

## 目标

帮助用户把云空间文件夹、知识库或文档库下的散乱文档组织成更清晰的目录结构。输出 `organize-plan.json`。

## 输入

- 必需：`inventory.json`
- 可选：用户给出的目标分类规则，例如“按项目”“按业务线”“按系统”“按年份”“按文档类型”
- 可选：outline 读取结果，用于按内容主题分类

## 分析规则

优先基于事实字段：

- 路径层级过深或过平。
- 空标题、重复标题、相似标题。
- 同一主题散落在多个目录。
- 云空间文件夹中在线文档和普通文件混放。
- 知识库 shortcut 复用或源文档散落。
- 文档标题包含“旧版”“废弃”“草稿”“临时”等治理信号。

模型推断必须写入 `reason` 和 `evidence`，不能把推断当事实。

## 生成计划

允许生成的 action：

| action | 说明 |
|-|-|
| `create_drive_folder` | 在 Drive 中创建目标文件夹 |
| `create_wiki_node` | 在知识库或文档库 / `my_library` 中创建目录节点 |
| `move_drive` | 移动 Drive 文件/文件夹 |
| `move_wiki_node` | 移动知识库节点 |
| `create_drive_shortcut` | 在 Drive 目标文件夹创建快捷方式 |
| `create_wiki_shortcut` | 在知识库中创建 shortcut 节点 |

禁止生成的 action：

- delete
- overwrite
- permission patch
- member remove
- owner transfer

每个 action 必须包含：

- source
- target
- reason
- evidence
- `requires_confirmation=true`
- dry-run command
- execute command

## 命令模板

Drive 创建文件夹：

```bash
lark-cli drive +create-folder --name "<folder_name>" --folder-token "<parent_folder_token>" --dry-run --as user
```

Drive 移动：

```bash
lark-cli drive +move --file-token "<token>" --type "<type>" --folder-token "<target_folder_token>" --dry-run --as user
```

Drive 快捷方式：

```bash
lark-cli drive +create-shortcut --file-token "<token>" --type "<type>" --folder-token "<target_folder_token>" --dry-run --as user
```

Wiki 创建节点：

```bash
lark-cli wiki +node-create --space-id "<space_id_or_my_library>" --parent-node-token "<parent_node_token>" --title "<title>" --obj-type docx --dry-run --as user
```

Wiki 移动节点：

```bash
lark-cli wiki +move --node-token "<node_token>" --target-parent-token "<target_parent_node_token>" --dry-run --as user
```

Wiki shortcut：

```bash
lark-cli wiki +node-create --space-id "<space_id>" --parent-node-token "<parent_node_token>" --node-type shortcut --origin-node-token "<source_node_token>" --title "<title>" --dry-run --as user
```

## 执行与验收

只有用户明确确认某个 `organize-plan.json` 后，才逐条执行。执行后写 `execution-log.json`，并重新跑 inventory 验证目标目录结构。

聊天回复要给出：

- 计划文件路径
- action 数量
- 需要确认的高影响动作
- 被阻塞的动作和原因
- 下一步确认方式
