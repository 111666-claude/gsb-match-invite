# match-invite

对局邀请：创建去重、有效期、人数上限与接受幂等，只用 Go 标准库。

```
go run ./cmd/match-invite --sample=dup
go run ./cmd/match-invite --sample=ttl
go run ./cmd/match-invite --sample=cap
go run ./cmd/match-invite --sample=idem
go run ./cmd/match-invite --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **创建去重**：同一个邀请码只创建一次，重复创建返回 false 且不改动原邀请。
- **有效期**：邀请存在时间达到 `ttlMs` 后失效，不能再被接受，`Pending` 也不再返回。
- **人数上限**：房间人数达到 `maxPlayers` 后 `Accept` 返回 false。
- **同人一次**：同一个玩家接受同一个邀请只成功一次，重复接受返回 false。
- **幂等**：同一个 `reqID` 只生效一次。
- **不变量**：房间人数不超过 `maxPlayers`；失效邀请不出现也不可接受；
  同一串创建与接受序列重复执行结果相同。
- **代价**：`Pending` 与 `Accept` 不许遍历全部邀请，`Scanned` 不随邀请数线性增长。
- **规模**：10 万个邀请、每秒 10 万次查询，单次摊还 O(1)。

## 输出契约（不改格式）

```
second=..
accepted=..
third=.. members=..
second=.. members=..
scanned=..
```

## 目录

```
invite.go              去重、有效期、上限与幂等
invite_test.go         单元测试
cmd/match-invite/      命令行入口
```
