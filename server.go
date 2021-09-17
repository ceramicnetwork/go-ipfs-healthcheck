package main

import (
  "fmt"
  "log"
  "net/http"
  "regexp"
)

type Status struct {
  Message string
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
  status := &Status{Message: "Health check succeeded"}
  fmt.Fprintf(w, status.Message)
}

var strictPath = regexp.MustCompile("^/$")

func createHandler(fn func (http.ResponseWriter, *http.Request), strict bool) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if strict {
      m := strictPath.FindStringSubmatch(r.URL.Path)
      if m == nil {
        http.NotFound(w, r)
        return
      }
    }
    fn(w, r)
  }
}

func main() {
  http.HandleFunc("/", createHandler(healthcheck, true))
  log.Fatal(http.ListenAndServe(":8080", nil))
}

