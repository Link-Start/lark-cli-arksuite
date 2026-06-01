# Wiki token routing

> **前置条件：** 先阅读 [`../SKILL.md`](../SKILL.md) 了解认证、全局参数、身份选择和错误处理。

Wiki 链接里的 `/wiki/<token>` 是知识库节点 token，不是底层文档、表格、Base、幻灯片或文件 token。处理 Wiki URL 时先解包，再按底层对象类型路由到对应 skill 或 API。

## 推荐入口

优先使用 `drive +inspect` 自动识别 URL 类型并解包 Wiki 节点：

```bash
lark-cli drive +inspect --url 'https://example.feishu.cn/wiki/<wiki_token>'
```

返回里的 `type`、`token`、`title`、`url` 是后续操作的 canonical 信息。后续命令使用返回的 canonical token，不要继续把 Wiki 节点 token 当作文件 token 传入。

## 手动解包

当 shortcut 不满足需求时，用 Wiki 节点接口查询底层对象：

```bash
lark-cli wiki spaces get_node --params '{"token":"<wiki_token>"}'
```

从返回结果中读取：

- `node.obj_type`：底层对象类型
- `node.obj_token`：底层对象 token
- `node.space_id`：知识空间 ID
- `node.title`：节点标题

## 路由表

| `obj_type` | 后续处理 |
|------------|----------|
| `docx` / `doc` | 文档正文读取和编辑走 `lark-doc`；评论、权限、导出、文件元数据走 `lark-drive` |
| `sheet` | 单元格、工作表、导出等表格操作走 `lark-sheets`；评论、权限、文件元数据走 `lark-drive` |
| `bitable` | 表、字段、记录、视图、表单、仪表盘、工作流走 `lark-base`；评论、权限、文件元数据走 `lark-drive` |
| `slides` | 幻灯片内容编辑走 `lark-slides`；导出、评论、权限、文件元数据走 `lark-drive` |
| `file` | 普通文件上传、下载、移动、删除、权限、评论走 `lark-drive` |
| `mindnote` | 按当前 CLI 支持能力优先走 `lark-drive`；缺口再评估原生 OpenAPI |

## 规则

- 不要把 Wiki 节点 token 直接当作 `doc_token`、`file_token`、`spreadsheet_token`、`app_token` 或 `presentation_token`。
- Wiki 节点、空间目录、成员管理等知识库结构操作保留 Wiki 节点 token，走 `lark-wiki`。
- 内容、评论、权限、导出等底层对象操作使用解包后的 `obj_token`。
- 如果用户只给 Wiki URL，先解析类型；不要凭标题、路径或上下文猜底层对象类型。
