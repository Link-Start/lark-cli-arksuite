# 权限治理 Command Patterns

本文只提供 `permission_governance` workflow 的具体 `lark-cli` 命令样例。只有进入对应 state 且需要拼装命令时才读取本文；命令可用范围仍以 [`lark-drive-workflow-permission-governance.md`](lark-drive-workflow-permission-governance.md) 的 `Command Map` 为准。

## 目录

- `目标解析`
- `目标发现`
- `事实读取`
- `写前确认与执行`

## 目标解析

```bash
lark-cli drive +inspect --url '<url>' --as user --format json
```

`/wiki/space/<space_id>` URL 是 Wiki space 范围，不要用 `drive +inspect` 当作单文档解析；直接提取 `space_id` 后进入 `DISCOVER_TARGETS`。

## 目标发现

发现 Wiki space / node 下目标：

```bash
lark-cli wiki +node-list \
  --space-id '<space_id>' --page-size 50 \
  --as user --format json

lark-cli wiki +node-list \
  --space-id '<space_id>' --parent-node-token '<node_token>' --page-size 50 \
  --as user --format json
```

发现 Drive folder 下目标：

```bash
lark-cli drive files list \
  --params '{"folder_token":"<folder_token>","page_size":200}' \
  --as user --format json

lark-cli drive files list \
  --params '{"folder_token":"<folder_token>","page_size":200,"page_token":"<PAGE_TOKEN>"}' \
  --as user --format json
```

## 事实读取

读取 metadata：

```bash
lark-cli drive metas batch_query \
  --data '{"request_docs":[{"doc_token":"<token>","doc_type":"<type>"}],"with_url":true}' \
  --as user --format json
```

读取 public permission：

```bash
lark-cli drive permission.public get \
  --params '{"token":"<token>","type":"<type>"}' \
  --as user --format json
```

按需读取访问统计：

```bash
lark-cli drive file.statistics get \
  --params '{"file_token":"<token>","file_type":"<type>"}' \
  --as user --format json
```

按需读取最近访问记录：

```bash
lark-cli drive file.view_records list \
  --params '{"file_token":"<token>","file_type":"<type>","page_size":50}' \
  --as user --format json
```

## 写前确认与执行

patch 前检查 manage-public permission：

```bash
lark-cli drive permission.members auth \
  --params '{"token":"<token>","type":"<type>","action":"manage_public"}' \
  --as user --format json
```

显式确认后 patch public permission：

```bash
lark-cli drive permission.public patch \
  --params '{"token":"<token>","type":"<type>"}' \
  --data '{"link_share_entity":"closed","external_access":false}' \
  --as user --yes --format json
```

显式确认后申请访问权限：

```bash
lark-cli drive +apply-permission \
  --token '<url>' \
  --perm view --remark '<reason>' --as user --format json

lark-cli drive +apply-permission \
  --token '<bare-token>' --type '<type>' \
  --perm view --remark '<reason>' --as user --format json
```

显式确认后更新 secure label：

```bash
lark-cli drive +secure-label-update \
  --token '<url>' \
  --label-id '<label-id>' --as user --format json

lark-cli drive +secure-label-update \
  --token '<bare-token>' --type '<type>' \
  --label-id '<label-id>' --as user --format json
```
