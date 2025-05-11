# Secure Proxy Management

## Requirements

- Go 1.20 or later

## Build

```bash
go install github.com/qiuzhanghua/autotag@latest
```

```bash
go mod tidy
```

```bash
go build -o bin/sproxy-mgmt
```

## Usage

### Set the API Key

create a file named `.env` in the root directory of the project and add the following lines to it:

```text
REDIS_URL=redis://localhost:6379/0
SECURE_PROXY_MGMT_PORT=9999
```

you can merge `.env` of `sproxy` and `sproxy-mgmt` into one file, if you put them in the same directory.

### Start the Proxy Server

```bash
bin/sproxy-mgmt
```

### call the API

use [httpie](https://httpie.io/) to call the API

```bash
http POST http://localhost:9999/api/add/alice/-1
# This will add a new key for alice with a TTL of -1 (no expiration).
# or
http POST http://localhost:9999/api/add/alice/forever
# or
http POST http://localhost:9999/api/add/alice/infinite

# 从当前开始，增加时分秒
http POST http://localhost:9999/api/add/bob/36h45m3s
# Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
```

### How to use the API

Use Java or any other language to call the API.
call add API to add a new key with a TTL.
and save the key in your database.

### API Status

- [x] `add` with TTL
- [ ] `list` by username
- [ ] `load` from database
- [ ] `remove` keys owned by username
- [ ] `delete` by key

## References

- [Redis URL format](https://pkg.go.dev/github.com/redis/go-redis/v9#ParseURL)

- [time.ParseDuration](https://pkg.go.dev/time#ParseDuration)

### Redis Cli Cmd

```redis
# 查看所有的key
keys *

# 查看某个key的值
get 95618478-99c6-4f20-935b-1792e400a65a
# "bob"

# 查看ttl
ttl 95618478-99c6-4f20-935b-1792e400a65a
# (integer) 160929

# 将有限的ttl改为无限
persist 95618478-99c6-4f20-935b-1792e400a65a

# 将无限的ttl改为有限
expire 95618478-99c6-4f20-935b-1792e400a65a 1000
# 1000秒

# 删除key
del 7746a258-9096-4b23-a344-8f44fb8a05c1 82f2a862-28ac-4094-b1d9-440e7d692a69

# 删除所有的key(！！小心使用！！)
flushdb

```
