package main

import (
  "flag"
  "fmt"
  "net/http"
  "strconv"
  "os"
  "github.com/fzzy/radix/redis"
  "time"
)

var listenPort int
var redisHost string
var redisPort string

func errHndlr(err error) {
  if err != nil {
    fmt.Println("error:", err)
    os.Exit(1)
  }
}

func init() {
  flag.IntVar(&listenPort, "port", 8080, "listen port")
  flag.StringVar(&redisHost, "redis_host", "127.0.0.1", "redis host")
  flag.StringVar(&redisPort, "redis_port", "6379", "redis port")
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

func redirect(w http.ResponseWriter, r *http.Request) {
  // fmt.Fprintf(os.Stderr,"Req: %s %s\n", r.Host, r.URL.Path)
  r_rule := rlookup(r.URL.Path)
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
  http.Redirect(w, r, r_target, r_code_int)
}

func main() {
  http.HandleFunc("/", redirect)
  fmt.Println("Listening on port:", strconv.Itoa(listenPort))
  fmt.Println("Using Redis host", redisHost)
  fmt.Println("Using Redis port", redisPort)
  http.ListenAndServe(":" + strconv.Itoa(listenPort), nil)
}
