提供docker service的启动，通过client启动而不是docker命令，因为agent在的容器不会有docker cli，只会把docker.sock挂载进来

启动不同类型的的服务，每个服务除了配置参数还可能存在配置文件

对于redis来说，可以参考的为：https://github.com/bitnami/containers/tree/main/bitnami/redis

```bash
docker pull bitnami/redis:6.2

docker run -e ALLOW_EMPTY_PASSWORD=yes -v /path/to/redis-persistence:/bitnami/redis/data bitnami/redis:latest
```

对于mysql来说，可以参考：https://github.com/bitnami/containers/tree/main/bitnami/mysql

```bash
docker pull bitnami/mysql:8.0

MYSQL_ROOT_PASSWORD=password123
```