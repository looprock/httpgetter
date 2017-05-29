package main

import (
     "fmt"
     "net/http"
     "time"
     "io/ioutil"
     "encoding/json"
		 "strings"
     "flag"
		 "github.com/riemann/riemann-go-client"
)

func MakeRequest(url string, ch chan<-string, riemannserver *string) {
  riemann := riemanngo.NewTcpClient(*riemannserver)
  err := riemann.Connect(5)
  if err != nil {
      panic(err)
  }
  start := time.Now()
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
  resp, err := client.Get(url)
  if err != nil {
      fmt.Println(err)
      return
  }
	a := strings.Split(url, "://")
	b := strings.Split(a[1], "?")
  c := strings.Replace(b[0], ".", "_", -1)
	d := strings.Replace(c, "/", "--", -1)
  secs := time.Since(start).Seconds()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      fmt.Println(err)
      return
  }
  tags := make([]string, 0)
  tags = append(tags, "http-check")
  riemanngo.SendEvent(riemann, &riemanngo.Event{
    Host: "http-check",
    Service: d,
    Metric: secs,
    Tags: tags,
  })
  defer riemann.Close()
  ch <- fmt.Sprintf("%s - site: %s, elapsed: %.2f, length: %d", time.Now().Format(time.RFC850), url, secs, len(body))
}

func main() {
  checkconfig := flag.String("config", "/etc/checks.json", "JSON site list file")
  riemannserver := flag.String("riemann", "127.0.0.1:5555", "riemann server and port")
  flag.Parse()
  data, err := ioutil.ReadFile(*checkconfig)
  if err != nil {
      fmt.Println(err)
      return
  }
  var slice []string
  err = json.Unmarshal(data, &slice)
  if err != nil {
      fmt.Println(err)
      return
  }
	for {
    ch := make(chan string)
	  for _,url := range slice {
	      go MakeRequest(url, ch, riemannserver)
	  }
	  // for range slice{
	  //   fmt.Println(<-ch)
	  // }
		time.Sleep(5 * time.Second)
	}
}
