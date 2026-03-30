
# redis docker run

- docker run

```golang
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v /Users/xxxxx/data/redis/redis.conf:/usr/local/etc/redis/redis.conf \
  -v /Users/xxxxx/data/redis/data:/data \
  --ulimit nofile=65535:65535 \
  --sysctl net.core.somaxconn=65535 \
  redis:7 \
  redis-server /usr/local/etc/redis/redis.conf
```


- redis.conf
```golang
# 绑定
bind 0.0.0.0
protected-mode no

# 端口
port 6379

# ===== ACL认证 =====
user default off
user admin on >password123! ~* +@all

# 内存策略
maxmemory 4gb
maxmemory-policy allkeys-lru

# 持久化
appendonly no
save ""

# 网络优化
tcp-backlog 511
timeout 0
tcp-keepalive 300

# 多线程
io-threads 4
io-threads-do-reads yes

# 日志
loglevel notice
```