# fanout-go

Send message for multiple subscribers

## Installing

```sh
go get -u github.com/xprgv/fanout-go
```

## Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/xprgv/fanout-go"
)

func main() {
    f := fanout.New[string]()

    sub1 := make(chan string, 1)
    f.AddSub(sub1)
    go func() {
        for {
            data, ok := <-sub1
            if !ok {
                fmt.Println("sub1 closed")
                return
            }
            fmt.Println("sub1:", data)
        }
    }()

    sub2 := make(chan string, 1)
    f.AddSub(sub2)
    go func() {
        for {
            data, ok := <-sub2
            if !ok {
                fmt.Println("sub2 closed")
                return
            }
            fmt.Println("sub2:", data)
        }
    }()

    for i := 0; i < 10; i++ {
        time.Sleep(time.Second)
        f.Publish(fmt.Sprint("iteration: ", i))
    }

    f.Close()

    time.Sleep(time.Second)
}
```
