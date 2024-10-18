/*

Golang L7 Dstat Source.
Code by t.me/uzerpanel

*/

package main

import (
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "strconv"
  "time"
)

type requests struct {
  Total_requests       int
  RequestsPerSecond    int
  MaxRequestsPerSecond int
  Date                 string
}

var rps requests
var index string

func main() {
  fmt.Println("(dstat by @uzerpanel | @tcpuzer) hacker dstat started on []:" + os.Args[1])
  go func() {
    for {
      fmt.Printf("RPS: " + strconv.Itoa(rps.RequestsPerSecond) + " | MAX RPS: " + strconv.Itoa(rps.MaxRequestsPerSecond) + " | TOTAL REQUESTS: " + strconv.Itoa(rps.Total_requests) + "\r")
      if rps.MaxRequestsPerSecond <= rps.RequestsPerSecond {
        rps.MaxRequestsPerSecond = rps.RequestsPerSecond
      }
      rps.RequestsPerSecond = 0
      rps.Date = time.Now().Format("01/02/2006 15:04:05")
      time.Sleep(1 * time.Second)
    }
  }()
  initHTTP()
}

func initHTTP() {
  index, _ = readFile("index.html")
  mux := http.NewServeMux()
  mux.HandleFunc("/", Index)
  mux.HandleFunc("/dstat", Dstat)
  mux.HandleFunc("/attack", Attack)

  err := http.ListenAndServe(":"+os.Args[1], mux)
  if errors.Is(err, http.ErrServerClosed) {
    fmt.Println("server closed")
  } else if err != nil {
    fmt.Println("cant start web server on: 127.0.0.1:" + os.Args[1])
    os.Exit(1)
  }
}

func Attack(w http.ResponseWriter, r *http.Request) {
  rps.RequestsPerSecond++
  rps.Total_requests++
  w.WriteHeader(666)
  io.WriteString(w, "attack me")
}

func Index(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, index)
}

func Dstat(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
  w.Header().Add("Content-Type", "application/json")
  w.Header().Add("Set-Cookie", "hacker=1")
  io.WriteString(w,
    "Content-Type: application/json\r\n"+
      "Set-Cookie: hacker=1\r\n"+
      "<title>@uzerpanel | PRIVATE DSTAT</title>"+
      `{"rps": `+strconv.Itoa(rps.RequestsPerSecond)+`, "max_rps": `+strconv.Itoa(rps.MaxRequestsPerSecond)+`, "total_requests": `+strconv.Itoa(rps.Total_requests)+`, "status": "online", "timestamp": "`+rps.Date+`"}`)
}

func readFile(filename string) (string, error) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return "", err
  }

  return string(data), nil
}