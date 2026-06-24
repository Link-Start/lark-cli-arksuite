# contact
> skill: lark-contact

## user_profiles batch_query
已知一批用户 id，要批量拿他们的个人状态（personal_status）和个性签名（description）

### Tips
- 默认不返回 personal_status 和 description，需在 query_option 里把 include_personal_status / include_description 显式置 true 才会带出对应字段
- user_ids 里 id 的类型要与 --user-id-type 一致（默认 open_id，可选 user_id / union_id）

### Examples

**批量查询个人状态和个性签名**
```bash
lark-cli contact user_profiles batch_query --data '{"user_ids":["ou_3a8b****6a7b"],"query_option":{"include_personal_status":true,"include_description":true}}'
```
