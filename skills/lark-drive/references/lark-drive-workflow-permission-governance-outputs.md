# 权限治理输出模板

本文只提供 `permission_governance` workflow 的用户可见输出模板。默认先给简短摘要；只有用户要求完整表格、需要写入确认，或结果大到需要结构化展示时才读取本文。

## 目录

- `输出策略`
- `用户语言映射`
- `定位与治理动作`
- `审计摘要`
- `容器安全诊断报告摘要`
- `可操作风险清单`
- `治理选择交互`
- `权限设置清单`
- `访问复核清单`
- `整改 dry-run`
- `批量权限申请确认`
- `确认请求`
- `最终摘要`

## 输出策略

- 单目标默认输出审计摘要。
- 容器目标默认输出安全诊断报告摘要：一句话结论、覆盖情况、风险分级、优先处理对象、建议下一步和剩余限制。
- 容器目标不要把风险按数量机械排序；外部公开、允许对外分享、缺失密级标签优先于复制 / 下载 / 评论这类依赖策略的候选项。
- 用户没有提供明确 policy 时，使用“候选风险 / 待复核 / 待策略确认”，不要写“违规 / 已泄露 / 已外部访问”。
- 风险对象展示按规模渐进披露：1-10 个全部展示；11-30 个展示全部高优先级待处理对象，中 / 低优先级只做分组摘要；31-100 个按高优先级分组展示 Top 5 和数量；100+ 个只展示分组统计和 Top 样例。
- 当摘要未展示全部风险对象时，必须明确“完整清单包含 <count> 条”，并提供生成 Markdown / CSV / 飞书文档风险清单或整改 dry-run 的下一步。
- 只要发现需要处理的对象，最终回复必须给出可执行下一步 CTA。不能因为默认只读，就只报告风险后结束。
- 完整风险清单是后续治理选择的输入；Markdown / CSV / 飞书文档报告必须使用同一套字段和稳定 `risk_id`。
- 写入前必须使用确认模板；权限申请、public-permission patch、secure-label update 分别确认。
- 最终回复必须包含已完成事项、验证结果和剩余限制；异步 owner 审批不能表述为已完成授权。

## 用户语言映射

面向用户的主结论优先使用业务语言；底层字段名只在证据或完整清单中保留。

| 底层字段 / 值 | 用户可见说法 |
|---------------|--------------|
| `link_share_entity=anyone_readable/anyone_editable` | 互联网公开链接候选风险 |
| `link_share_entity=tenant_readable` | 公司内知道链接可读 |
| `link_share_entity=tenant_editable` | 公司内知道链接可编辑 |
| `link_share_entity=closed` | 未开启链接分享 |
| `external_access=true` | 允许分享到组织外；不等于已经存在外部协作者 |
| `external_access=false` | 不允许分享到组织外 |
| `share_entity=anyone` | 较多人可添加或管理协作者 |
| `share_entity=same_tenant` | 公司内成员可添加或管理协作者 |
| `share_entity=only_full_access` | 仅有管理权限的人可管理协作者 |
| `security_entity` is not `only_full_access` | 复制 / 下载 / 打印范围需要按策略复核 |
| `comment_entity=anyone_can_view` | 可查看者都能评论 |
| `sec_label_name` missing | 缺少密级标签 |

## 定位与治理动作

风险对象必须能让用户直接定位和处理：

- 摘要中的每个优先处理对象必须包含 `path/title`、`URL`、`type`、风险原因、关键证据和建议动作。
- 完整清单、访问复核清单、整改 dry-run 和写入确认都必须包含 URL。缺少 URL 时，展示 token / node_token，并说明 URL 未能获取。
- 同名文档、shortcut 或副本必须用 path + URL 区分；不要只输出 title。
- 完整风险清单中的每条记录必须有稳定 `risk_id`，格式为 `PG-001`、`PG-002`。`risk_id` 在同一次诊断和后续 dry-run / 确认 / 验证中保持不变。
- 建议动作必须和风险类型绑定：互联网公开链接优先建议关闭链接分享或收紧为组织内；允许对外分享优先建议 owner 复核或关闭对外分享；缺少密级标签优先建议补齐密级；复制 / 下载 / 评论范围只在用户 policy 明确时建议收紧。
- 写入动作只能作为下一步选项或确认请求出现。不要在诊断摘要里暗示已经执行缩权。

## 摘要清单展开规则

容器安全诊断的摘要必须兼顾可读性和可治理性。不要用固定 Top N 代替可处理清单。

| 风险对象数 | 摘要默认展示 | 必须提供的下一步 |
|------------|--------------|------------------|
| `0` | 只展示覆盖情况、未覆盖能力和剩余限制 | 如需更细审计，可生成权限设置清单 |
| `1-10` | 展示全部风险对象 | 可直接按 `risk_id` 生成 dry-run 或写入确认 |
| `11-30` | 展示全部高优先级待处理对象；中 / 低优先级做分组摘要 | 生成完整风险清单 artifact，或按风险分组生成 dry-run |
| `31-100` | 每个高优先级风险分组展示 Top 5，附未展示数量 | 生成 Markdown / CSV / 飞书文档完整风险清单 |
| `100+` | 只展示分组统计、Top 样例和覆盖限制，不内联长表 | 强烈建议生成结构化风险清单后再选择治理范围 |

高优先级待处理对象包括：互联网公开链接、允许对外分享、允许对外分享且缺少 / 低于 policy 密级标签、公司内可编辑链接、协作者管理范围较宽。复制 / 下载 / 打印、评论范围在用户未提供明确 policy 时归入“待策略确认”，不要挤占高优先级清单。

摘要中的每个待处理对象必须包含 `risk_id`、path/title、URL、type、风险原因、关键证据和建议动作。对同一底层文档的多个 Wiki 入口或 shortcut，必须用 URL 区分；如果建议合并治理，在建议动作里说明它们指向同一底层对象。

## 审计摘要

```text
目标：<title> (<type>)
URL：<url-or-token-if-url-unavailable>
结论：<合规 / 待确认风险 / 无法完整判断>
证据：
- link_share_entity=<value>
- external_access=<value>
- share_entity=<value>
- security_entity=<value>
- comment_entity=<value>
- sec_label_name=<value-or-missing>
限制：<unsupported_checks or none>
建议动作：<read-only next step or proposed remediation>
```

## 容器安全诊断报告摘要

```text
已完成只读安全诊断，没有做任何权限修改。

一句话结论：<未发现互联网公开链接 / 存在互联网公开链接候选风险>；<external_access_count> 个文档允许对外分享，<missing_label_count> 个文档缺少密级标签。建议优先复核 <top_priority_group_or_paths>。

覆盖情况：
- 当前身份可见目标：<visible_count>
- 已成功检查公开权限：<permission_checked_count>
- 读取失败 / 已删除 / 无权限：<failed_count>
- 未覆盖能力：<collaborator_list / inheritance / audit_log / view_records / none>

风险分级：
- 高优先级：<internet_public_count> 个互联网公开链接候选；<external_access_count> 个允许对外分享；其中 <external_without_label_count> 个同时缺少密级标签。
- 中优先级：<tenant_link_count> 个公司内知道链接可访问 / 可编辑；<wide_share_count> 个协作者管理范围较宽。
- 待策略确认：<security_count> 个复制 / 下载 / 打印范围待复核；<comment_count> 个评论范围待复核。
- 无法判断：<unsupported_or_unverified_summary>。

高优先级待处理清单：
> 按 `摘要清单展开规则` 展示。每个对象必须包含 `risk_id` 和 URL；缺少 URL 时展示 token / node_token 和原因。若没有高优先级对象，只展示中优先级或待策略确认分组摘要。

- <risk_id> <path-or-title> (<type>)
  URL: <url-or-token-if-url-unavailable>
  风险：<why high priority>
  证据：<short field evidence>
  建议动作：<recommended action>

未完全展开：
- 完整风险清单包含 <risk_manifest_count> 条；本摘要已展示 <shown_count> 条，未展示 <hidden_count> 条。
- 未展示分组：<risk_group=count summary or none>

建议下一步：
- 生成完整风险清单 artifact，包含 `risk_id`、URL、证据字段、建议动作和 `selected` 列。
- 基于 risk_id、风险分组、owner、路径、URL 或 artifact 中 `selected=true` 的行生成只读整改 dry-run。
- 只针对最高优先级目标进入写入确认流程，例如关闭互联网公开链接或收紧对外分享；写入前仍需二次确认。
- 按 owner / 密级生成复核清单。
- 继续读取访问记录，判断低活跃高暴露。

剩余限制：
- <do not claim collaborator-list verification if unsupported>
- <external_access=true only means sharing outside is allowed, not that external collaborators exist>
- <missing view_records / DLP / AI index status / audit log limitations>
```

## 可操作风险清单

完整风险清单用于让用户选择后续治理范围。Markdown / CSV / 飞书文档报告都必须包含以下字段；如果某种格式无法完整展示嵌套证据，使用短文本摘要，保留 `risk_id` 和 URL。

```text
范围：<wiki_space / wiki_node / drive_folder> <name-or-id>
生成时间：<timestamp>
用途：用户可按 risk_id、risk_group、owner、path、URL 或 selected=true 选择治理对象。

| risk_id | Path | URL | Type | Owner | risk_group | evidence | recommended_action | current_setting | target_setting | selected | decision | status | skip_reason |
|---------|------|-----|------|-------|------------|----------|--------------------|-----------------|----------------|----------|----------|--------|-------------|
| PG-001 | <path> | <url-or-token> | <type> | <owner-or-unknown> | <risk_group> | <short evidence> | <recommended-action> | <field=value> | <field=value-or-owner-review> | false | undecided | pending | <none-or-reason> |
```

字段规则：

- `risk_id` 按风险优先级和 path 稳定排序生成；同一次诊断中不得重复。
- `selected` 默认 `false`；用户可在 CSV / 飞书文档表格中改为 `true`，或在聊天中直接说 “处理 PG-001、PG-003”。
- `decision` 表示用户决策：`undecided`、`keep`、`dry_run`、`confirm_write`、`skip`。
- `status` 表示执行状态：`pending`、`dry_run_ready`、`confirmed`、`executed`、`verified`、`failed`、`skipped`。
- `target_setting` 是建议目标状态，不代表已执行；没有明确 policy 时只能写 owner review / policy review。

## 治理选择交互

用户基于完整风险清单继续治理时，Agent 必须先解析选择范围，再生成只读 dry-run：

```text
可接受的用户选择：
- 处理 PG-001、PG-003、PG-008，把互联网公开链接关闭。
- 先处理所有 risk_group=internet_public_link，不处理 external_access_only。
- 把 CSV / 飞书文档里 selected=true 的行生成整改 dry-run。
- PG-003 先跳过，只处理 PG-001。

Agent 必须回复：
- 已选择对象数：<count>
- 选择来源：<risk_id list / risk_group / selected=true / URL / path>
- 将执行的下一步：生成 dry-run；不执行写入
- 需要跳过或重新确认的对象：<missing risk_id / unsupported / changed_since_report / no manage_public>
```

如果用户选择来自旧报告或外部 artifact，生成 dry-run 前必须对所选目标重新读取当前权限。当前设置和报告快照不一致时，标记为 `changed_since_report`，不要直接沿用旧字段执行。

## 权限设置清单

```text
范围：<wiki_space / wiki_node / drive_folder> <name-or-id>

| Path | URL | Type | link_share_entity | external_access | share_entity | security_entity | comment_entity | sec_label_name | 建议动作 | 限制 |
|------|-----|------|-------------------|-----------------|--------------|-----------------|----------------|----------------|----------|------|
| <path> | <url-or-token> | <type> | <value> | <value> | <value> | <value> | <value> | <value-or-missing> | <recommended-action> | <unsupported-or-none> |
```

## 访问复核清单

```text
范围：<wiki_space / wiki_node / drive_folder / explicit_list> <name-or-id>
复核对象数：<count>

| Owner | Path | URL | Type | 风险标签 | 当前权限摘要 | 最近访问证据 | 建议动作 |
|-------|------|-----|------|----------|--------------|--------------|----------|
| <owner-or-unknown> | <path> | <url-or-token> | <type> | <labels> | <link/external/share/security/comment> | <uv/pv/last_view_or_unknown> | <keep / tighten / owner review / unsupported> |

限制：<unsupported_checks / discovery_blockers / none>
```

## 整改 dry-run

```text
将生成整改计划，不执行写入：
- 范围：<scope>
- 选择来源：<risk_id list / risk_group / selected=true artifact / URL list>
- 候选目标数：<count>
- 计划执行命令：<command family>
- 重新读取：已对所选目标重新读取当前权限；changed_since_report=<count>
- 字段变更：
  - <risk_id> <path> (<url-or-token>): <field> <old> -> <new>
- 跳过项：<unsupported / no manage_public / unsupported type / missing policy>
- 验证方式：执行后重新读取 <metadata/public_permission>
- 有限回滚范围：<public_permission_snapshots fields or not applicable>

请确认是否进入写入确认。
```

## 批量权限申请确认

```text
将逐个发起 <view / edit> 权限申请：
- 候选目标数：<count>
- 命令类型：drive +apply-permission
- 风险：write；每个请求都会通知 owner
- 执行方式：按候选列表顺序逐个调用，失败项会单独记录

候选示例：
- <risk_id> <title> (<type>, <url-or-token>)：<reason>

请确认是否对上述候选目标发起权限申请。
```

## 确认请求

```text
将执行 <operation>：
- 目标：<risk_id> <title> (<type>, <url-or-token>)
- 命令类型：<command family>
- 风险：<risk_level>
- 字段变更：
  - <field>: <old> -> <new>
- 验证方式：执行后重新读取 <metadata/public_permission>
- 有限回滚材料：<public_permission_snapshots or not applicable>

请确认是否执行。
```

## 最终摘要

```text
已完成：<read checks / writes>
验证：<fresh read result or async owner approval note>
清单状态：<risk_id status updates / not applicable>
回滚材料：<public_permission_snapshots / not applicable>
剩余限制：<unsupported_checks / partial facts / approvals>
```
