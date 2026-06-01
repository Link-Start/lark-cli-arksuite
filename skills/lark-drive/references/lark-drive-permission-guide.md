# drive permission guide

> **前置条件：** 先阅读 [`../SKILL.md`](../SKILL.md) 了解 Drive 入口、身份路由和 Wiki token 解包规则；再阅读 [`../../lark-shared/SKILL.md`](../../lark-shared/SKILL.md) 了解认证与全局参数。

## 权限申请

- 向文档 owner 申请 view / edit 权限时，优先使用 `drive +apply-permission`。
- `drive +apply-permission` 仅支持 user 身份；不要用 bot 反复重试用户个人资源。
- 如果用户只给 token，先尽量用 `drive +inspect` 或上下文恢复 URL / 类型，再申请权限。

## 公开权限错误码

调用 `drive.permission.public.patch` 修改公开权限时，常见业务错误码按下表处理：

| 错误码 | 处理建议 |
|--------|----------|
| `91009` | 当前文件已设置更高安全等级或受组织策略限制，无法开放到目标范围；提示用户在文档 UI 或安全设置中处理 |
| `91010` | 文件 owner 或管理员设置不允许修改公开权限；需要 owner / 管理员调整策略 |
| `91011` | 当前身份没有修改公开权限的权限；切换到有权限的 user，或让 owner 授权 |
| `91012` | 目标公开范围不符合当前租户策略；不要继续扩大权限，提示用户按组织策略处理 |

如果错误信息与上表不一致，以接口返回的 `msg` / `error` 为准，并给出下一步可执行动作。

## 授权当前应用访问文档

当需要把文档权限授予**当前应用（bot）自身**时，先获取当前 bot 的 `open_id`，再授权给该 `open_id`：

```bash
# 1. 获取当前应用的 open_id
lark-cli api GET /open-apis/bot/v3/info --as bot

# 2. 授权当前应用访问文档
lark-cli drive permission.members create \
  --params '{"token":"<doc_token>","type":"<resource_type>"}' \
  --data '{"member_type":"openid","member_id":"<bot_open_id>","perm":"view","type":"user"}'
```

说明：

- 此方式只适用于授权给当前应用自身。
- 授权给其他用户时，直接使用对方的 `open_id`，不要调用 bot info。
- 不要转移 owner，除非用户明确要求且已说明影响。
- `<resource_type>` 常见值：`doc`、`docx`、`sheet`、`bitable`、`file`、`folder`、`wiki`、`slides`。

## 身份判断

- 用户个人云空间、用户拥有的文档、公开权限申请、密级标签等场景优先 `--as user`。
- bot 创建或已授权给 bot 的资源可使用 `--as bot`。
- bot 因资源不可见失败时，说明资源可见性问题；建议切换 user 身份或先给 bot 授权，不要盲目重试。
