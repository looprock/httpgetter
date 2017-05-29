package main

import (
     "fmt"
     "os"
     "net/http"
     "time"
     "io/ioutil"
		 "strings"
		 //"github.com/riemann/riemann-go-client"
)

// this is for documentation of what riemanngo.Event is expecting
type Revent struct {
	Ttl         float32
	// Time        time.Time
	Tags        []string
	// Host        string // Defaults to os.Hostname()
	State       string
	Service     string
	Metric      interface{} // Could be Int, Float32, Float64
	Description string
	// Attributes  map[string]string
}

func printSlice(s []Revent) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

//func MakeRequest(url string, ch chan<-riemanngo.Event) {
func MakeRequest(url string, ch chan<-Revent) {
  start := time.Now()
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
  resp, _ := client.Get(url)
  // service = cmd_args[0].split('://')[1].replace('.', '_').replace('/', '--').split("?")[0]
	a := strings.Split(url, "://")
	b := strings.Split(a[1], "?")
  c := strings.Replace(b[0], ".", "_", -1)
	d := strings.Replace(c, "/", "--", -1)
  secs := time.Since(start).Seconds()
  ioutil.ReadAll(resp.Body)
	// var tags [1]string
  tags := make([]string, 0)
  tags = append(tags, "fakesrv")
	//tags[0] = "fakesrv"
	//t1 := riemanngo.Event{d, secs, tags}
  t1 := Revent{60, tags, "OK", d, secs, "http check" }
  ch <- t1
}

func main() {
	// c := riemanngo.NewTcpClient("127.0.0.1:5555")
	// err := c.Connect(5)
	// if err != nil {
	//     panic(err)
	// }
	// https://play.golang.org/p/ireemFw2Xi
	for {
	  start := time.Now()
		//ch := make(chan riemanngo.Event)
    ch := make(chan Revent)
	  for _,url := range os.Args[1:]{
	      go MakeRequest(url, ch)
	  }

    // events = []riemanngo.Event {
    //     riemanngo.Event{
    //         Service: "hello",
    //         Metric:  100,
    //         Tags: []string{"hello"},
    //     },
    // riemanngo.Event{
    //         Service: "goodbye",
    //         Metric:  200,
    //         Tags: []string{"goodbye"},
    //     },
    // }

    // events = []riemanngo.Event {}

    // probably need something like this:
    // https://stackoverflow.com/questions/19818878/slice-append-from-channels
    events := []Revent {}
	  for range os.Args[1:]{
	    fmt.Println(<-ch)
      // var chanevent Revent
      //chanevent <-ch
      //events = append(events, chanevent)
	  }
    fmt.Print(len(events))
	  fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
		fmt.Printf("sleeping 5\n")
		time.Sleep(5 * time.Second)
	}
	// defer c.Close()
}
