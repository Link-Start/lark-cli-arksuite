# 知识资产整理 Artifact 协议

> 用途：让 inventory、organize、permission-audit 等 workflow 可串联、可恢复、可评测。字段保持最小稳定；缺失信息用空值、`warnings` 或 `unsupported_checks` 表达，不要编造。

## 目录约定

```text
./lark-drive-knowledge/<run-id>/
  scope.json
  inventory.json
  organize-plan.json
  permission-audit.json
  execution-log.json
  report.md
```

`run-id` 建议使用 `YYYYMMDD-HHMMSS-<short-scope>`，例如 `20260526-143000-wiki-space`.

## scope.json

记录用户目标和本次实际处理范围。

```json
{
  "run_id": "",
  "requested_by_user": "",
  "scope_type": "drive_folder|wiki_space|wiki_node|my_library|mixed",
  "root": {
    "url": "",
    "space_id": "",
    "folder_token": "",
    "node_token": ""
  },
  "limits": {
    "max_depth": -1,
    "page_limit": 0,
    "content_read": "none|outline|targeted"
  },
  "generated_at": ""
}
```

## inventory.json

记录事实清单。所有后续分析优先消费它。

```json
{
  "run_id": "",
  "summary": {
    "total": 0,
    "by_source": {},
    "by_type": {},
    "warnings_count": 0
  },
  "items": [
    {
      "source": "drive|wiki|my_library",
      "title": "",
      "path": "",
      "url": "",
      "token": "",
      "type": "folder|docx|doc|sheet|bitable|file|slides|mindnote|wiki|shortcut",
      "space_id": "",
      "node_token": "",
      "obj_token": "",
      "obj_type": "",
      "folder_token": "",
      "parent_token": "",
      "depth": 0,
      "has_child": false,
      "owner": "",
      "created_time": "",
      "modified_time": "",
      "evidence": []
    }
  ],
  "warnings": [],
  "generated_at": ""
}
```

## organize-plan.json

只表达计划，不代表已经执行。

```json
{
  "run_id": "",
  "mode": "plan",
  "summary": {
    "actions_count": 0,
    "requires_confirmation_count": 0
  },
  "actions": [
    {
      "id": "act-001",
      "action": "create_drive_folder|create_wiki_node|move_drive|move_wiki_node|create_drive_shortcut|create_wiki_shortcut",
      "source": {},
      "target": {},
      "reason": "",
      "evidence": [],
      "risk": "read|write|high-risk-write",
      "requires_confirmation": true,
      "dry_run_command": "",
      "execute_command": ""
    }
  ],
  "blocked": [],
  "warnings": [],
  "generated_at": ""
}
```

第一版不生成 delete、overwrite、permission patch、owner transfer 动作。

## permission-audit.json

记录权限审计事实、推断和能力边界。

```json
{
  "run_id": "",
  "summary": {
    "audited_items": 0,
    "risk_findings": 0,
    "unsupported_checks": []
  },
  "items": [
    {
      "title": "",
      "path": "",
      "url": "",
      "token": "",
      "type": "",
      "wiki_space_members": [],
      "public_permission": {},
      "risk_findings": [
        {
          "type": "",
          "severity": "low|medium|high",
          "fact": "",
          "inference": "",
          "suggestion": "",
          "evidence": [],
          "requires_confirmation": true
        }
      ],
      "unsupported_checks": []
    }
  ],
  "warnings": [],
  "generated_at": ""
}
```

## execution-log.json

仅在用户明确确认执行 organize plan 后生成。

```json
{
  "run_id": "",
  "plan_file": "",
  "results": [
    {
      "action_id": "act-001",
      "status": "success|failed|skipped",
      "command": "",
      "result": {},
      "error": ""
    }
  ],
  "generated_at": ""
}
```
