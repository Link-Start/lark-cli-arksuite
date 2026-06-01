
> **术语说明：** 飞书云空间也常被称为"云盘"或"云存储"，三者指的是同一个产品，是飞书官方的云端文件存储与管理中心。

## 关键规则速览

- Wiki 链接里的 `/wiki/<token>` 不是底层文件 token；优先用 `drive +inspect` 解包，完整规则见 [`lark-wiki-token-routing.md`](../lark-shared/references/lark-wiki-token-routing.md)。
- 本地文件导入为在线文档统一走 `drive +import`：Excel / CSV / `.base` -> `--type bitable`，Office / Markdown / HTML / TXT -> `--type docx`，PPTX -> `--type slides`。
- 原生 `.md` 文件内容读写、patch、diff 走 [`lark-markdown`](../lark-markdown/SKILL.md)；把 Markdown 转成在线 docx 才走 `drive +import --type docx`。
- 评论默认只查未解决评论；review / 审阅场景优先局部评论。评论细节见 [`lark-drive-comments-guide.md`](references/lark-drive-comments-guide.md)。
- 权限申请、安全标签、搜索等用户个人资源场景优先 `--as user`；bot 只能处理自己可见或已授权的资源。

## 快速决策

- 用户要**搜文档 / Wiki / 电子表格 / 多维表格 / 云空间对象**，优先使用 `lark-cli drive +search`。自然语言里"最近我编辑过的"、"我创建的"（`--mine`，owner 语义）、"最近一周我打开过的 xxx"、"某人 owner 的 docx" 等直接映射到扁平 flag。
- 用户要把本地 `.xlsx` / `.csv` / `.base` 导入成 Base / 多维表格 / bitable，第一步必须使用 `lark-cli drive +import --type bitable`。
- 用户要把本地 `.md` / `.docx` / `.doc` / `.txt` / `.html` 导入成在线文档，使用 `lark-cli drive +import --type docx`。
- 用户要把本地 `.xlsx` / `.xls` / `.csv` 导入成电子表格，使用 `lark-cli drive +import --type sheet`。
- 用户要把本地 `.pptx` 导入成飞书幻灯片，使用 `lark-cli drive +import --type slides`；当前 PPTX 导入上限是 500MB。
- 用户要在 Drive 里上传、创建、读取、局部 patch 或覆盖更新**原生 `.md` 文件**（不是导入成 docx），切到 [`lark-markdown`](../lark-markdown/SKILL.md)。
- 用户要比较原生 `.md` 文件的历史版本差异，或比较远端 Markdown 与本地草稿，切到 [`lark-markdown`](../lark-markdown/SKILL.md) 的 `lark-cli markdown +diff`；需要版本号时先用 `drive +version-history`。
- 用户要查看、下载、回滚或删除文件的历史版本，使用 `drive +version-history`、`drive +version-get`、`drive +version-revert`、`drive +version-delete`；自动化场景优先 `--as bot`。
- 用户要在云空间里新建文件夹，优先使用 `lark-cli drive +create-folder`。
- 用户要把本地文件上传到知识库 / 文档库里的某个 wiki 节点下，仍然使用 `lark-cli drive +upload --wiki-token <wiki_token>`；不要误切到 `wiki` 域命令。
- 用户要修改标题，可用 `drive files patch` 传 `new_title`，支持 docx、sheet、bitable、file、wiki、folder 类型。

## 身份路由

| 场景 | 身份建议 |
|------|----------|
| 搜索用户可见资源、申请权限、查看/更新安全标签 | 使用 `--as user`；这些 shortcut 是 user-only 或强依赖用户可见范围 |
| 用户个人云空间、用户拥有的文档/文件夹 | 默认 `--as user` |
| bot 自己创建的资源、已授权给 bot 的资源、自动化版本操作 | 可用 `--as bot`；版本命令自动化场景优先 bot |
| bot 因资源不可见失败 | 不要反复重试；提示用户切 `--as user`，或先把资源授权给当前应用 |

## 核心概念

- 直接文档 URL（如 `/docx/`、`/doc/`、`/sheets/`、`/drive/folder/`）通常可从路径直接取得对应 token。
- Wiki URL（`/wiki/<token>`）必须先解析到底层 `obj_type` 和 `obj_token`，再决定后续调用哪个域。
- `drive +inspect` 是跨类型 URL 检视的首选入口；当它输出 `type` 和 `token` 后，后续命令使用该 canonical token。
- 原生 API 调用前先运行 `lark-cli schema drive.<resource>.<method>` 查看 `--params` / `--data` 结构；不要猜字段。

## 评论与权限

- 添加评论优先使用 [`drive +add-comment`](references/lark-drive-add-comment.md)；查询、统计、回复限制和 reaction 规则见 [`lark-drive-comments-guide.md`](references/lark-drive-comments-guide.md) 与 [`lark-drive-reactions.md`](references/lark-drive-reactions.md)。
- 权限申请优先使用 [`drive +apply-permission`](references/lark-drive-apply-permission.md)；公开权限错误码、授权当前应用访问文档等规则见 [`lark-drive-permission-guide.md`](references/lark-drive-permission-guide.md)。

## 不在本 skill 范围

- 文档正文读取、改写、追加、替换、图片/附件插入：切到 [`lark-doc`](../lark-doc/SKILL.md)。
- 电子表格单元格、工作表、筛选、公式等表内操作：切到 [`lark-sheets`](../lark-sheets/SKILL.md)。
- Base 表、字段、记录、视图、表单、仪表盘、工作流等表内操作：切到 [`lark-base`](../lark-base/SKILL.md)。
- 知识空间、Wiki 节点、空间成员管理：切到 [`lark-wiki`](../lark-wiki/SKILL.md)。
- Drive 原生 Markdown 文件的创建、读取、patch、overwrite、diff：切到 [`lark-markdown`](../lark-markdown/SKILL.md)。

## 参考

- [`lark-wiki-token-routing.md`](../lark-shared/references/lark-wiki-token-routing.md) — Wiki token / obj_token / obj_type 路由规则
- [`lark-drive-comments-guide.md`](references/lark-drive-comments-guide.md) — 评论查询、统计、回复和 review 落点规则
- [`lark-drive-permission-guide.md`](references/lark-drive-permission-guide.md) — 权限错误码和应用授权规则
