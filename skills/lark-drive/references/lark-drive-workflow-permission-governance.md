# lark-drive 权限治理 Workflow

Workflow id: `permission_governance`

Risk / Structure: `R2` / `S2`

本文实现已注册的权限治理 workflow。执行前必须先读取 [`lark-drive-workflow.md`](lark-drive-workflow.md)，并遵循其中的共享执行协议、Artifact Contract 和 Workflow Loading 规则。

## 必读上下文

执行本 workflow 前，先阅读 [`../../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) 了解认证、全局参数和安全规则。

## 适用范围

当用户要求检查或治理 Drive / Docs / Wiki 资产访问权限时，使用本 workflow：

- "检查这个文档是不是对外公开"
- "帮我收紧这个文档权限"
- "给这些文档做权限风险排查"
- "检查这个知识库下所有文档是否存在安全风险"
- "查看这个文件夹下所有文档的权限设置情况"
- "生成 owner 访问复核清单，找出低活跃但权限过宽的文档"
- "找出这个知识库下哪些文档我没有权限，并帮我申请查看 / 编辑权限"
- "评估这些文档是否适合进入 AI Agent / RAG 检索上下文"
- "申请查看 / 编辑这个文档"
- "调整这个文档的密级标签"
- "看看谁可以分享、复制、下载、评论"

目标可以是一个明确 URL / token、小规模明确列表、Wiki space / Wiki node，或 Drive folder。遇到 Wiki space / Drive folder 等容器范围时，先执行只读 `DISCOVER_TARGETS` 并产出覆盖摘要；这里的"所有文档"只表示当前身份在确认范围内可枚举到的文档。任何写入都必须再次询问并获得确认。

## 非目标

本 workflow 不处理：

- 目录组织、迁移、归档或清理；这类需求应使用知识整理 workflow。
- 内容审查、过期内容判断或知识质量评分。
- owner transfer、协作者创建 / 撤销、成员列表审计；本 workflow 只输出 owner 复核候选，不执行 owner 转移。
- 文件夹自身公开权限审计或修复。`drive permission.public get` / `patch` 不支持 `type=folder`；必须记录到 `unsupported_checks`，然后继续读取文件夹下其他支持的文档事实。
- 当前身份无法枚举到的不可见文档的完整发现；只能处理已发现目标，或用户显式提供的 URL / token。
- 未按范围确认的批量写入。

不要声称已完成协作者列表验证：当前 CLI surface 没有 `permission.members list` shortcut。

## Progressive Load Map

本表只规定每个 state 需要加载的额外上下文；命令可用范围以 `Command Map` 为准。未进入对应 state 前，不要预读无关 reference。

| State | Required Reference |
|-------|--------------------|
| `PARSE_INTENT` | 本文件、[`lark-drive-workflow.md`](lark-drive-workflow.md)、[`../../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) |
| `TARGET_INSPECT` | [`lark-drive-inspect.md`](lark-drive-inspect.md)；需要具体命令样例时读取 [`lark-drive-workflow-permission-governance-commands.md`](lark-drive-workflow-permission-governance-commands.md) |
| `DISCOVER_TARGETS` | 容器范围时读取 [`../../lark-wiki/references/lark-wiki-node-list.md`](../../lark-wiki/references/lark-wiki-node-list.md) 或 [`lark-drive-files-list.md`](lark-drive-files-list.md)；需要具体命令样例时读取 [`lark-drive-workflow-permission-governance-commands.md`](lark-drive-workflow-permission-governance-commands.md) |
| `FACT_READ` | `lark-cli schema drive.metas.batch_query`；涉及公开权限时再读取 `lark-cli schema drive.permission.public.get`；涉及活跃度、访问复核或生命周期判断时再读取 `lark-cli schema drive.file.statistics.get` 和 `lark-cli schema drive.file.view_records.list`；需要具体命令样例时读取 [`lark-drive-workflow-permission-governance-commands.md`](lark-drive-workflow-permission-governance-commands.md) |
| `RISK_ASSESS` | 本文件的 `Risk Classification` |
| `EXEC_CONFIRM` | 只为用户选择的动作读取 [`lark-drive-apply-permission.md`](lark-drive-apply-permission.md)、[`lark-drive-secure-label.md`](lark-drive-secure-label.md)，或 `lark-cli schema drive.permission.public.patch`；需要确认模板时读取 [`lark-drive-workflow-permission-governance-outputs.md`](lark-drive-workflow-permission-governance-outputs.md) |
| `EXECUTE` | 复用 `EXEC_CONFIRM` 已加载且已确认的写命令上下文；需要具体命令样例时读取 [`lark-drive-workflow-permission-governance-commands.md`](lark-drive-workflow-permission-governance-commands.md) |
| `VERIFY` | 复用 `FACT_READ` 阶段使用的 read schemas |

## Runtime State Extension

本 workflow 在共享 `Artifact Contract` 基础上扩展以下字段：

| Field | Meaning |
|-------|---------|
| `intent` | `audit`、`list_permission_settings`、`access_review`、`remediation_dry_run`、`tighten_public_permission`、`apply_permission` 或 `secure_label_update` |
| `target_scope` | `single_resource`、`explicit_list`、`wiki_space`、`wiki_node` 或 `drive_folder`，包含用户原始输入和已确认范围 |
| `targets` | 标准化直接目标和发现目标列表，包含原始输入、解析后的 type、token、title、URL、path，以及存在时的 wiki node / object data |
| `discovered_targets` | 从 Wiki space / Wiki node / Drive folder 只读发现出的可审计文档目标 |
| `discovery_blockers` | 发现阶段因权限、分页、API 覆盖、目标类型或工具预算导致的未覆盖范围 |
| `coverage_summary` | 容器范围的发现数量、已审计数量、unsupported 数量和 partial 原因 |
| `metadata_facts` | `drive metas batch_query` 结果，包括 title、owner、URL，以及返回时的 `sec_label_name` |
| `public_permission_facts` | 对支持目标执行 `drive permission.public get` 的结果 |
| `activity_facts` | 用户要求活跃度、最近访问、闲置暴露或访问复核时，读取 `drive file.statistics get` / `drive file.view_records list` 得到的访问证据 |
| `manage_public_auth` | patch 写入前，以 `action=manage_public` 执行 `permission.members auth` 的结果 |
| `risk_findings` | 基于证据的发现和置信度；每项必须包含稳定 `risk_id`、path、title、type、URL 或可替代定位信息、风险标签、关键证据和建议动作 |
| `risk_manifest` | 完整风险清单 artifact 的行数据；包含 `risk_id`、定位信息、风险分组、当前设置、建议目标设置、`selected`、`decision`、`status` 和 `skip_reason` |
| `selected_risk_items` | 用户通过 risk_id、风险分组、owner、路径、URL、CSV `selected=true` 或飞书文档表格选择出的待处理目标 |
| `access_review_items` | 面向 owner / 项目负责人的复核清单，包含 path、URL、owner、风险标签、当前权限设置、最近访问证据和建议动作 |
| `permission_request_candidates` | 权限申请候选目标；只来自已发现且可构造申请请求的目标 |
| `remediation_plan` | dry-run 或已确认整改计划，包含目标、字段 diff、跳过原因、验证方式和执行顺序 |
| `public_permission_snapshots` | public-permission patch 前保存的原始字段快照，仅用于说明有限回滚范围 |

## Execution State Machine

| State | Protocol Step | Agent MUST Do | User-Facing Output | wait_for_user | Next State |
|-------|---------------|---------------|--------------------|---------------|------------|
| `PARSE_INTENT` | `route` / `scope` | 解析 intent、target scope、desired policy，以及只读审计、权限申请还是修复模式 | 范围确认；如果缺少目标或期望动作，只问一个澄清问题 | 缺少 target / action，或容器范围需要用户确认时为 `true` | `TARGET_INSPECT` |
| `TARGET_INSPECT` | `scope` | 解析单资源、明确列表、Wiki space / node、Drive folder；保留原始 URL、scope type、canonical token/type | 目标范围表，包含 scope、title/type/token status | 除非解析失败，否则为 `false` | `DISCOVER_TARGETS` or `FACT_READ` |
| `DISCOVER_TARGETS` | `scope` / `read` | 对 Wiki space / node 或 Drive folder 递归只读枚举，归一化为 `discovered_targets`；记录 `discovery_blockers` | 发现进度和覆盖摘要；不展示内部 cursor/token，除非用户要求 | 除非发现范围无法确认或全部被阻断，否则为 `false` | `FACT_READ` |
| `FACT_READ` | `read` | 对直接目标或 `discovered_targets` 执行 `drive metas batch_query`；对支持的非 folder 目标执行 `drive permission.public get`；在用户要求活跃度 / 访问复核 / 生命周期判断时读取访问统计和访问记录 | 权限事实摘要、coverage summary、activity facts 和 unsupported checks | 除非所有目标都被 auth 阻断，否则为 `false` | `RISK_ASSESS` |
| `RISK_ASSESS` | `assess/plan` | 对证据分类；如用户提供 policy，则对照 policy；构建可定位风险清单、访问复核清单、dry-run 整改计划或候选修复计划；完整清单必须生成稳定 `risk_id` | 带 URL 和 risk_id 的 findings、confidence、review items、建议动作和下一步 CTA | `true` | `EXEC_CONFIRM` or `DONE` |
| `EXEC_CONFIRM` | `confirm` | 展示准确写入范围、command family、target count、risk、verification method | 确认请求 | `true` | `EXECUTE` or `DONE` |
| `EXECUTE` | `execute` | 只执行 `Command Map` 中已确认的写入 | 进度 / 结果摘要 | 除非被阻断，否则为 `false` | `VERIFY` |
| `VERIFY` | `verify` | 重新执行支持的读取，并与目标状态对比 | 验证表和剩余缺口 | `false` | `DONE` |
| `DONE` | `done` | 停止 | 最终回复，包含完成事项、验证结果和剩余风险 | `false` | End |

## Command Map

本 workflow 只能使用以下 command families：

| State | Allowed Command Families | Purpose |
|-------|--------------------------|---------|
| `TARGET_INSPECT` | `drive +inspect` | 解析 URL、type、canonical token、title 和 wiki unwrap data |
| `DISCOVER_TARGETS` | `wiki +node-list` | 递归发现 Wiki space / node 下当前身份可见的节点 |
| `DISCOVER_TARGETS` | `drive files list` | 递归发现 Drive folder 下当前身份可见的文件和子文件夹 |
| `FACT_READ` | `drive metas batch_query` | 读取 title、URL、owner 和 secure-label metadata |
| `FACT_READ` | `drive permission.public get` | 读取支持类型的 public/link/external/copy/comment/share settings |
| `FACT_READ` | `drive file.statistics get` | 在用户要求活跃度、闲置暴露、生命周期或访问复核时读取文件访问统计 |
| `FACT_READ` | `drive file.view_records list` | 在用户要求最近访问人、访问复核或低活跃证据时读取访问记录 |
| `EXEC_CONFIRM` | `drive +secure-label-list` | 提议 label update 前解析可用 secure-label IDs |
| `EXEC_CONFIRM` | `drive permission.members auth` | public-permission patch 前检查 `action=manage_public` |
| `EXECUTE` | `drive +apply-permission` | 向 owner 提交 view/edit access request；只允许单目标、小列表或已明确确认的候选列表逐个执行 |
| `EXECUTE` | `drive permission.public patch` | 修改已确认的 public/link settings；必须传 `--yes` |
| `EXECUTE` | `drive +secure-label-update` | 设置已确认的 secure-label ID |
| `VERIFY` | `drive metas batch_query`, `drive permission.public get` | 验证支持的 metadata 和 public-permission changes |

## Command Patterns

本入口不内联命令样例。需要拼装具体 `lark-cli` 命令时，按当前 state 读取 [`lark-drive-workflow-permission-governance-commands.md`](lark-drive-workflow-permission-governance-commands.md)。命令是否允许执行仍以 `Command Map` 和写入规则为准。

## Discovery Rules

容器范围只能先做只读发现和覆盖摘要，不能在发现阶段执行权限申请、权限 patch 或密级更新。

通用规则：

1. "所有文档"只表示当前身份在确认范围内可枚举到的文档。不可见、无权限、API 不返回或工具预算不足的部分必须进入 `discovery_blockers` 或 `unsupported_checks`。
2. 发现阶段必须生成稳定 `path`。不要只保存 title；同名文档必须能通过 path 或 token 区分。
3. 只把支持类型加入可审计目标：`doc`、`sheet`、`file`、`wiki`、`bitable`、`docx`、`mindnote`、`slides`。
4. `folder` 只作为递归容器，不执行 `permission.public get` / `patch`。`shortcut`、`catalog` 或缺少 stable token/type 的条目必须记录为 unsupported，除非后续 API 明确解析出支持目标。
5. 对大范围目标输出进度时，只展示已扫描容器数、已发现目标数、已审计目标数、剩余队列或 blocker；不要默认展示内部 page token / cursor。

Wiki space / node 发现：

1. `/wiki/space/<space_id>` 直接解析为 `target_scope=wiki_space`。不要因为 `drive +inspect` 对该 URL 返回 not found 就停止。
2. 用 `wiki +node-list --space-id <space_id>` 读取根节点；当节点 `has_child=true` 时，用该节点的 `node_token` 继续递归读取子节点。
3. Wiki 节点必须同时保留 `node_token`、`obj_token` 和 `obj_type`。权限读取优先用 `type=wiki` + `node_token` 表达 Wiki 节点权限；元数据补充可使用 `obj_type` + `obj_token`。
4. 如果节点只有 `obj_token` / `obj_type`，但无法确认 Wiki 节点权限 token，保留该目标为 partial，并在 `unsupported_checks` 中说明只能读取底层对象或无法完整判断 Wiki 节点权限。

Drive folder 发现：

1. `/drive/folder/<folder_token>` 解析为 `target_scope=drive_folder`。文件夹自身公开权限不支持；继续枚举其子文档。
2. 按 [`lark-drive-files-list.md`](lark-drive-files-list.md) 递归处理 `data.files`、`has_more` 和 `next_page_token`。不要把第一页数量当作完整范围。
3. 只对返回项中的 `folder` 继续递归；对子文档按 `type + token` 归一化为 `discovered_targets`。
4. 如果某个目录分页失败、无 continuation token、权限不足或 API 报错，只阻断该目录分支，并在 `discovery_blockers` 中记录；继续处理其他可枚举分支。

## Fact Read Rules

1. `drive metas batch_query` 单次最多 200 个 `request_docs`；当 `targets` 或 `discovered_targets` 超过 200 个时，必须分批读取并合并结果。
2. `drive permission.public get` 没有批量读取接口；对支持目标逐个读取。单个目标失败时记录 `unsupported_checks` 或 `partial`，不要阻断其他目标。
3. 对 Wiki 发现目标，公开权限读取优先使用 `type=wiki` + `node_token`；metadata 可使用 `obj_type` + `obj_token` 补充 title、owner、URL 和 `sec_label_name`。
4. 当 intent 是 `list_permission_settings` 时，只输出权限设置清单和覆盖限制，不主动生成修复计划。
5. `drive file.statistics get` 和 `drive file.view_records list` 只在用户要求最近访问、活跃度、闲置暴露、访问复核，或用户提供的 policy 明确依赖活跃度时执行；不要为普通权限审计默认读取访问记录。
6. 访问统计 / 访问记录当前只对 `doc`、`docx`、`sheet`、`bitable`、`mindnote`、`wiki`、`file` 作为支持类型处理。其他类型必须进入 `unsupported_checks`，不能推断活跃度。
7. `view_records` 是访问证据，不是权限列表。没有返回访问记录只能表述为“未获得最近访问证据”或“低活跃候选”，不能表述为“无人有权限”。

## Risk Classification

以下标签只能作为 evidence labels。除非用户提供明确 policy，否则不要把它们直接表述为绝对违规：

| Evidence | Suggested Label |
|----------|-----------------|
| `link_share_entity=anyone_readable` or `anyone_editable` | 外部链接暴露候选风险 |
| `link_share_entity=tenant_editable` | 租户范围可编辑候选风险 |
| `external_access=true` | 已启用外部分享 |
| `share_entity=anyone` | 广泛协作者管理候选风险 |
| `security_entity` is not `only_full_access` | 复制 / 下载 / 打印范围候选风险 |
| `comment_entity=anyone_can_view` | 广泛评论范围候选风险 |
| `sec_label_name` missing or lower than user-provided policy | 密级标签待复核候选风险 |
| `statistics.uv` / `statistics.pv` 很低或缺少最近访问证据，同时存在广泛可访问设置 | 低活跃高暴露候选风险 |
| `owner_id` / owner metadata 缺失、无法解析，或用户明确要求 owner 复核但 owner 不可用 | owner 待复核候选风险 |
| 广泛可访问设置 + 缺少 / 较低密级标签 + 用户要求 AI / Agent / RAG 前置治理 | AI 检索暴露候选风险 |

如果缺少 policy，必须把发现表述为“待确认风险”，并给出准确字段和值。

容器级安全诊断不要按命中数量机械排序。默认按以下用户决策优先级组织：

1. `link_share_entity=anyone_readable/anyone_editable`：最高优先级，表述为“互联网公开链接候选风险”。
2. `external_access=true`：高优先级，表述为“允许对外分享候选风险”；必须说明这不等于已经存在外部协作者。
3. `external_access=true` + `sec_label_name` missing / lower than policy：高优先级，建议优先补标签或 owner 复核。
4. `link_share_entity=tenant_readable/tenant_editable`、`share_entity=anyone`：中优先级，表述为“组织内广泛可见 / 可管理候选风险”。
5. `security_entity` / `comment_entity` 这类复制、下载、打印、评论范围：除非用户提供 policy，否则放入“待策略确认”，不要压过外部分享风险。
6. 无法读取协作者名单、继承链、DLP、AI 索引状态、访问记录时，放入“无法判断 / 未覆盖”，不要推断风险不存在。

`AI 检索暴露候选风险` 只是基于权限和标签的代理标签。除非另有工具明确返回索引状态，否则不要声称某个文档已经被 Agent、Copilot 或 RAG 索引。

## 写入规则

- Public-permission patch 属于高风险写入。请求确认前，必须展示 target title、token、current setting、desired setting 和准确 field changes。
- 如果 `manage_public_auth.auth_result=false`，禁止 patch。告诉用户需要具备 manage-public 权限的用户，或由 owner 操作。
- `drive permission.public get` 只用于 `drive +inspect` 或 `DISCOVER_TARGETS` 可解析且受支持的目标类型：`doc`、`sheet`、`file`、`wiki`、`bitable`、`docx`、`mindnote`、`slides`。
- 不要 patch 已解析类型不支持的字段。对于 wiki 目标，必须省略 schema 明确标注为 wiki 不支持的字段。
- 不要在同一个 confirmation 中合并 secure-label update 和 public-permission patch；必须分别确认。
- `drive +apply-permission` 默认不批量执行；每次调用都会向 owner 发送通知。
- 容器范围内的"统一申请权限"必须先产出 `permission_request_candidates`。未展示候选目标、数量、权限类型和 owner 通知影响前，禁止调用 `drive +apply-permission`。
- 用户显式确认批量权限申请后，也必须逐个目标顺序调用 `drive +apply-permission`，并在结果中区分已发起申请、失败、无法构造申请请求和未发现目标。
- 用户要求“生成整改方案 / dry-run / 先看看会改什么”时，只生成 `remediation_plan`，不执行任何写命令。dry-run 必须包含 target count、field changes、跳过原因、验证方式和有限回滚范围。
- 用户基于完整风险清单选择对象时，必须先解析 `risk_id`、风险分组、URL 或 artifact 中 `selected=true` 的行，生成 `selected_risk_items`。无法匹配到当前 `risk_manifest` 的选择必须要求用户重新确认或重新读取清单。
- 针对 `selected_risk_items` 生成 dry-run 前，必须重新读取所选目标的 `drive permission.public get`；如果当前设置和清单快照不同，标记为 `changed_since_report` 并跳过或要求用户确认更新后的计划。
- 执行 `drive permission.public patch` 前，必须把当前 `public_permission_facts` 中会被改动的字段保存为 `public_permission_snapshots`。该快照只用于 public-permission 字段的有限回滚说明，不覆盖协作者、owner、继承权限或密级标签。
- 如果用户要求批量收紧权限，必须按风险分层和目标顺序逐个执行；失败项进入结果清单，不要因为单个失败而重复执行已成功目标。
- 遇到 secure-label downgrade error `1063013` 时，停止重试，并告诉用户需要在文档 UI 中完成审批。

## 未来扩展边界

以下能力已有部分 CLI surface 或用户价值，但不要在当前 workflow 中作为可执行分支直接调用：

- `drive permission.members create` 可创建协作者权限，但当前 workflow 不做协作者 grant / update / revoke；未来需要单独定义授权对象解析、最小权限、确认模板和验证方式。
- `drive permission.members transfer_owner` 属于 owner transfer 高风险写入；当前 workflow 只输出 owner 复核候选，不执行 owner 转移。
- `wiki +member-list` 可作为 Wiki space 成员治理的读侧事实来源；当前 workflow 只治理文档 / 节点 / 文件夹下可发现文档的权限，不做 space member governance。
- 当前 CLI 没有 `permission.members list`、完整继承链、DLP 扫描、AI 索引状态、审计日志和跨平台权限事实。遇到这些需求必须记录为 `unsupported_checks` 或建议新增独立 workflow。

## 输出策略

- 默认 summary-first：单目标输出简短审计摘要；容器目标输出安全诊断报告摘要，而不是字段计数堆叠。
- 容器安全诊断必须包含一句话结论、覆盖情况、风险分级、优先处理对象、建议下一步和剩余限制。
- 优先处理对象必须可定位：默认展示 path/title、URL、type、风险原因、关键证据和建议动作。缺少 URL 时，必须展示 token、node_token 或无法生成 URL 的原因。
- 容器摘要不能固定只展示 Top N 样例。风险对象 1-10 个时全部展示；11-30 个时展示全部高优先级待处理对象；31 个以上按高优先级分组展示 Top 样例，并明确完整清单总数和生成 artifact 的下一步。
- 发现风险后不要只结束报告。必须给出下一步 CTA，例如查看完整风险清单、生成只读整改 dry-run、按 owner / 密级生成复核清单，或对高优先级目标进入写入确认流程。
- 完整风险清单是后续治理选择的输入，不是普通报告。每一行必须有稳定 `risk_id`，让用户可以按编号、风险分组、owner、路径、URL、CSV `selected=true` 或飞书文档表格行选择处理范围。
- 面向用户时优先使用业务语言：如“允许对外分享”“公司内知道链接可读”“谁可以复制 / 下载 / 打印”“谁可以评论”。底层字段名只作为证据补充，不作为主结论。
- 完整权限设置清单、访问复核清单、整改 dry-run、写入确认和最终摘要模板，按需读取 [`lark-drive-workflow-permission-governance-outputs.md`](lark-drive-workflow-permission-governance-outputs.md)。
- 不要默认创建文件、飞书文档或长表格；只有用户要求、结果过大，或 workflow 明确需要结构化确认时再生成。
- 最终回复必须包含已完成事项、验证结果和剩余限制。异步 owner 审批只能表述为“已发起申请”，不能表述为已完成授权。
