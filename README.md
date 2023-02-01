> for develop, your host must have docker

# Usage

## config

新增config

```bash
curl localhost:8000/swarm/config/create --form name=testname --form "content=@./test.conf"
```

列出所有config

```bash
curl localhost:8000/swarm/config/list
```

更新config

> 注意，不允许修改config内容: Error response from daemon: rpc error: code = InvalidArgument desc = only updates to Labels are allowed

```bash
curl localhost:8000/swarm/config/update --form name=testname --form "content=@./test.conf"
```

删除config

```bash
curl localhost:8000/swarm/config/remove -H "Content-Type: application/json" -d '{"name": "testname"}'
```

## common service

- start: 启动服务
- desc: 描述如何启动的服务

### redis

```bash
agentcli start redis
```