# drive comments guide

> **前置条件：** 先阅读 [`../SKILL.md`](../SKILL.md) 了解 Drive 入口、身份路由和 Wiki token 解包规则；再阅读 [`../../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) 了解认证与全局参数。

## 创建评论

- 添加评论优先使用 `drive +add-comment`。
- 全文评论：未传 `--block-id` 时默认启用，也可显式传 `--full-comment`；支持 `docx`、旧版 `doc` URL、白名单扩展名普通文件，以及最终解析为这些类型的 Wiki URL。
- 局部评论：传 `--block-id` 时启用；支持 `docx`、`sheet`、`slides`，以及最终解析为这些类型的 Wiki URL。docx block ID 可通过 `docs +fetch --api-version v2 --detail with-ids` 获取；sheet 使用 `<sheetId>!<cell>`；slides 使用 `<slide-block-type>!<xml-id>`。
- 如果 Wiki URL 解包后不是 `doc` / `docx` / `file` / `sheet` / `slides`，不要使用 `+add-comment`。
- Drive 普通文件评论仅支持平台允许的文件类型，常见后缀包括 `.md`、`.txt`、`.json`、`.csv`、`.go`、`.js`、`.py`、`.pptx`、`.png`、`.jpg`、`.jpeg`、`.zip`、`.mp3`、`.mp4`；普通文件只支持全文评论。
- `--content` 需要传 `reply_elements` JSON 数组字符串，例如 `--content '[{"type":"text","text":"正文"}]'`。
- 如果直接调用原生评论 V2 API，先执行 `lark-cli schema drive.file.comments.create_v2`；全文评论省略 `anchor`，局部评论传 `anchor.block_id`。

## Review 场景

- 代码、文案、方案 review 场景优先创建局部评论。
- 不同问题尽量拆成多条局部评论，便于逐条解决和追踪。
- Drive 普通文件不支持局部锚点，只能创建全文评论。

## 内容与转义

- 评论内容使用 `reply_elements` 表达，不要把纯 Markdown 当作结构化评论正文。
- `drive +add-comment` 会处理常见文本转义；直接调原生 API 时需要自行转义 `<` 和 `>`，避免被平台当作标签片段。
- 需要在评论中表达代码或 JSON 时，优先作为普通文本元素传入，不要猜测不在 schema 里的富文本字段。

## 查询与统计

- `drive file.comments list` 默认必须传 `is_solved:false`，即仅查询未解决评论。
- 即使用户说“所有评论”“全部评论”“把评论都列出来”，只要没有明确要求包含已解决评论，仍按未解决评论查询。
- 只有用户明确要求包含已解决评论时，才省略 `is_solved`。

```bash
# 默认查询：仅未解决评论
lark-cli drive file.comments list --params '{"file_token":"<DOC_TOKEN>","file_type":"docx","is_solved":false}'

# 用户明确要求包含已解决评论时
lark-cli drive file.comments list --params '{"file_token":"<DOC_TOKEN>","file_type":"docx"}'
```

- `items` 是评论卡片列表，不是平铺的互动消息列表。
- 创建第一条评论时会同时创建该卡片里的第一条 reply；真正承载正文的是 `item.reply_list.replies`。
- 统计“评论数”或“评论卡片数”：统计 `items` 数量，分页场景累加所有页。
- 统计“回复数”：所有 `item.reply_list.replies` 数量之和减去 `items` 数量。
- 统计“总互动数”：所有 `item.reply_list.replies` 数量之和，包含每张评论卡片的首条评论。
- 如果 `item.has_more=true`，继续调用 `drive file.comment.replys list` 拉全该卡片下的回复，再做全量统计。

## 排序

- 只有用户明确提到“最新评论”“最后评论”“最早评论”时，才需要按 `create_time` 排序。
- 排序前必须先处理分页，拉完所有评论；不要只取第一页排序。
- “最新评论” / “最后评论”：按 `create_time` 降序取第一条。
- “最早评论”：按 `create_time` 升序取第一条。
- 用户只说“第一条评论”时，直接使用 `drive file.comments list` 返回的第一条，不额外排序。

## 回复限制

- 全文评论不支持回复：`is_whole=true` 的评论无法添加回复。
- 已解决评论不支持回复：`is_solved=true` 的评论无法添加回复。
- 用户要回复某条评论但该评论不能回复时，只提示限制原因；不要自动替用户寻找其他可回复评论。

## 批量查询与 reaction

- `drive file.comments batch_query` 只用于已知评论 ID 后的批量查询，需要传入具体 comment ID 列表。
- 遍历、统计、获取最新/最早评论等场景使用 `drive file.comments list`。
- reaction（表情、点赞、谁点了什么、添加/删除表情）场景先阅读 [`lark-drive-reactions.md`](lark-drive-reactions.md)。
