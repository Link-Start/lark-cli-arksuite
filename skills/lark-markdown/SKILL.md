---
name: lark-markdown
version: 1.2.0
description: "飞书 Markdown：查看、创建、上传、编辑和比较 Drive 原生 Markdown 文件。当用户需要创建或编辑 `.md` 文件、读取、修改、局部 patch 或比较差异时使用。不负责：把 Markdown 导入成飞书在线文档 docx（走 lark-drive 的 drive +import --type docx）、文件搜索/权限/评论/移动/删除等云空间管理操作（走 lark-drive）。"
metadata:
  requires:
    bins: ["lark-cli"]
  cliHelp: "lark-cli markdown --help"
---

# markdown (v1)

**CRITICAL — 开始前 MUST 先用 Read 工具读取 [`../lark-shared/SKILL.md`](../lark-shared/SKILL.md)，其中包含认证、权限处理**

## 快速决策

- 用户要**上传、创建一个原生 `.md` 文件**，使用 `lark-cli markdown +create`
- 用户要**比较原生 `.md` 文件的历史版本差异**，或比较远端 Markdown 与本地草稿，使用 `lark-cli markdown +diff`
- 用户要**读取 Drive 里某个 `.md` 文件内容**，使用 `lark-cli markdown +fetch`
- 用户要对 Markdown 文件做**局部文本替换 / 正则替换**，优先使用 `lark-cli markdown +patch`
- 用户要**覆盖更新 Drive 里某个 `.md` 文件内容**，使用 `lark-cli markdown +overwrite`
- 用户要先拿 Markdown 文件的历史版本号，再做比较/下载/回滚，先用 [`lark-drive`](../lark-drive/SKILL.md) 的 `lark-cli drive +version-history`
- 用户要把本地 Markdown **导入成在线新版文档（docx）**，不要用本 skill，改用 [`lark-drive`](../lark-drive/SKILL.md) 的 `lark-cli drive +import --type docx`
- 用户要对 Markdown 文件做**rename / move / delete / 搜索 / 权限 / 评论**等云空间（云盘/云存储）操作，不要留在本 skill，切到 [`lark-drive`](../lark-drive/SKILL.md)

## 身份策略

- 所有 markdown shortcut 同时支持 `--as user` 和 `--as bot`
- 访问用户个人云空间或用户拥有的文件时，优先使用 `--as user`
- CI / 自动化场景、bot 自己创建或已被授权访问的 Markdown 文件，可使用 `--as bot`
- 如果 bot 因看不到用户资源而失败，不要反复重试；提示用户改用 `--as user` 或先给 bot 授权

## 核心边界

- 本 skill 处理的是 **Drive 中作为普通文件存储的 Markdown**，不是 docx 文档
- `--name` 和本地 `--file` 文件名都必须显式带 `.md` 后缀；不满足时 shortcut 会直接报错
- `--content` 支持：
  - 直接传字符串
  - `@file` 从本地文件读取内容
  - `-` 从 stdin 读取内容
- `markdown +patch` 的内部语义是：**先完整下载 Markdown，再本地替换，再整文件覆盖上传**
- `markdown +patch` 不是服务端原子 patch；它是 CLI 侧编排出来的局部更新能力
- `markdown +patch` 当前只支持**单组** `--pattern` / `--content`
- `markdown +patch` 替换后的最终内容**不能为空**；Drive 不支持 0 字节 Markdown，且空结果通常代表误清空，CLI 会直接报错，不会上传空文件
- `--file` 只接受本地 `.md` 文件路径

## Shortcuts（推荐优先使用）

Shortcut 是对常用操作的高级封装（`lark-cli markdown +<verb> [flags]`）。有 Shortcut 的操作优先使用。

| Shortcut | 说明 |
|----------|------|
| [`+create`](references/lark-markdown-create.md) | 在 Drive 中创建原生 Markdown 文件 |
| [`+diff`](references/lark-markdown-diff.md) | 比较两个远端 Markdown 历史版本，或比较远端 Markdown 与本地文件 |
| [`+fetch`](references/lark-markdown-fetch.md) | 读取 Drive 中的原生 Markdown 文件 |
| [`+patch`](references/lark-markdown-patch.md) | 通过下载、替换、覆盖上传的方式局部更新 Markdown 文件 |
| [`+overwrite`](references/lark-markdown-overwrite.md) | 覆盖更新 Drive 中已有的 Markdown 文件 |

## 参考

- [lark-shared](../lark-shared/SKILL.md) — 认证和全局参数
- [lark-drive](../lark-drive/SKILL.md) — Drive 文件管理、导入 docx、move/delete/search 等
