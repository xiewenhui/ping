package ping

import (
  "io/ioutil"
  "log"
  "runtime"
  "strings"
  "sync"
  "testing"
  "time"
  "os"
)

func TestPing(t *testing.T) {
  if os.Geteuid() != 0 {
    log.Fatalln("Need Root Privilege To Run")
  }

  log.Printf("cpunum: %v \n", runtime.NumCPU())
  runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores

  filename := "/tmp/ips"
  bytes, err := ioutil.ReadFile(filename)

  if err != nil {
    t.Fatalf("filename not exist: %v", filename)
  }

  ips := strings.Split(string(bytes[:]), ",")

  var wg sync.WaitGroup

  cocurrent := 2
  timeOutMs := 300
  tokens := make(chan int, cocurrent)

  start := time.Now().Unix()
  log.Println(ips)
  log.Println(len(ips))

  for i := 0; i < len(ips); i++ {
    tokens <- 1
    wg.Add(1)

    ip := strings.Trim(ips[i], "\n")
    go func(h string) {
      defer wg.Done()
      alive := Ping(h, timeOutMs)
      log.Printf("[%v]:[%v]\n", h, alive)
      // time.Sleep(3000 * time.Millisecond)
      <- tokens
    }(ip)
  }
  wg.Wait()

  end := time.Now().Unix()
  log.Println(start, end)
}

