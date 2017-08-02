# POC: GORRed
redirect server in GO, backed by Redis

## Usage
```
Usage of ./gorred_osx:
  -port int
        listen port (default 8080)
  -redis_host string
        redis host (default "127.0.0.1")
  -redis_port string
        redis port (default "6379")
```

## Expectations
Service expects to lookup redirects via HGETALL
```
  redis-cli hgetall /becroak_test/
    1) "target"
    2) "/becroak_redirest/"
    3) "code"
    4) "302"
```

you can populate your redis instance with included `redirects.csv`
```
  export REDIS_HOST=127.0.0.1
  for i in $(cat redirects.csv); do  \
    url=$(echo $i | cut -d, -f1);    \
    code=$(echo $i | cut -d, -f2);   \
    target=$(echo $i | cut -d, -f3); \
    redis-cli -h $REDIS_HOST hset $url target $target; \
    redis-cli -h $REDIS_HOST hset $url code $code; \
  done
```

