package main

import (
  "flag"
  "fmt"
  "github.com/valyala/fasthttp"
  "strconv"
  "os"
  "github.com/fzzy/radix/redis"
  "time"
)

var listenAddr string
var redisHost string
var redisPort string
var compress bool

func errHndlr(err error) {
  if err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }
}

func init() {
  flag.StringVar(&listenAddr, "listen", ":9000", "listen address")
  flag.StringVar(&redisHost, "redis_host", "127.0.0.1", "redis host")
  flag.StringVar(&redisPort, "redis_port", "6379", "redis port")
  flag.BoolVar(&compress, "compress", false, "whether to enable transparent response compression")
  flag.Parse()
}

func rlookup(k string) map[string]string {
  red, err := redis.DialTimeout("tcp", redisHost + ":" + redisPort, time.Duration(2)*time.Second)
  errHndlr(err)
  defer red.Close()

  r := red.Cmd("select", 0)
  errHndlr(r.Err)

  myhash, err := red.Cmd("hgetall", k).Hash()
  errHndlr(err)

  // switch {
  // case len(
  // r_target) == 0:
  //   fmt.Println("no redirect target here")
  // }
  red.Close()
  return myhash
}

func redirectHandler(ctx *fasthttp.RequestCtx) {
  r_path := string(ctx.Path())
  r_rule := rlookup(r_path)
  r_target := r_rule["target"]
  r_code := r_rule["code"]
  switch {
  case len(r_target) == 0:
    r_target = "/catchall_default"
  }
  switch {
  case len(r_code) == 0:
    r_code = "301"
  }
  r_code_int, err := strconv.Atoi(r_code)
  errHndlr(err)
  ctx.Redirect(r_target, r_code_int)
}

func main() {
  fmt.Println("Listening on address", listenAddr)
  fmt.Println("Using Redis host", redisHost)
  fmt.Println("Using Redis port", redisPort)
  fmt.Println("Compression enabled:", compress)

  h := redirectHandler
  if compress == true {
    h = fasthttp.CompressHandler(h)
  }

  if err := fasthttp.ListenAndServe(listenAddr, h); err != nil {
    fmt.Println("Error in ListenAndServe: %s", err)
  }
}
